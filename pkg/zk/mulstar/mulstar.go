package zkmulstar

import (
	"crypto/rand"

	"github.com/cronokirby/safenum"
	"github.com/taurusgroup/cmp-ecdsa/internal/hash"
	"github.com/taurusgroup/cmp-ecdsa/pkg/math/arith"
	"github.com/taurusgroup/cmp-ecdsa/pkg/math/curve"
	"github.com/taurusgroup/cmp-ecdsa/pkg/math/sample"
	"github.com/taurusgroup/cmp-ecdsa/pkg/paillier"
	"github.com/taurusgroup/cmp-ecdsa/pkg/pedersen"
)

type (
	Public struct {
		// C = Enc₀(?;?)
		C *paillier.Ciphertext

		// D = (x ⨀ C) ⨁ Enc₀(y;ρ)
		D *paillier.Ciphertext

		// X = gˣ
		X *curve.Point

		// Verifier = N₀
		Verifier *paillier.PublicKey
		Aux      *pedersen.Parameters
	}
	Private struct {
		// X ∈ ± 2ˡ
		X *safenum.Int

		// Rho = ρ = Nonce D
		Rho *safenum.Nat
	}
)

func (p Proof) IsValid(public Public) bool {
	if !arith.IsValidModN(public.Verifier.N(), p.W) {
		return false
	}
	if !public.Verifier.ValidateCiphertexts(p.A) {
		return false
	}
	if p.Bx.IsIdentity() {
		return false
	}
	return true
}

func NewProof(hash *hash.Hash, public Public, private Private) *Proof {
	N0Big := public.Verifier.N()
	N0 := safenum.ModulusFromNat(new(safenum.Nat).SetBig(N0Big, N0Big.BitLen()))

	verifier := public.Verifier

	alpha := sample.IntervalLEps(rand.Reader)

	r := sample.UnitModN(rand.Reader, N0)

	gamma := sample.IntervalLEpsN(rand.Reader)
	m := sample.IntervalLEpsN(rand.Reader)

	A := public.C.Clone().Mul(verifier, alpha)
	A.Randomize(verifier, r)

	commitment := &Commitment{
		A:  A,
		Bx: *curve.NewIdentityPoint().ScalarBaseMult(curve.NewScalarInt(alpha)),
		E:  public.Aux.Commit(alpha, gamma),
		S:  public.Aux.Commit(private.X, m),
	}

	e := challenge(hash, public, commitment)

	// z₁ = e•x+α
	z1 := new(safenum.Int).Mul(e, private.X, -1)
	z1.Add(z1, alpha, -1)
	// z₂ = e•m+γ
	z2 := new(safenum.Int).Mul(e, m, -1)
	z2.Add(z2, gamma, -1)
	// w = ρ^e•r mod N₀
	w := new(safenum.Nat).ExpI(private.Rho, e, N0)
	w.ModMul(w, r, N0)

	return &Proof{
		Commitment: commitment,
		Z1:         z1.Big(),
		Z2:         z2.Big(),
		W:          w.Big(),
	}
}

func (p *Proof) Verify(hash *hash.Hash, public Public) bool {
	if !p.IsValid(public) {
		return false
	}

	verifier := public.Verifier

	if !arith.IsInIntervalLEps(p.Z1) {
		return false
	}

	e := challenge(hash, public, p.Commitment)

	if !public.Aux.Verify(p.Z1, p.Z2, p.E, p.S, e.Big()) {
		return false
	}

	z1 := new(safenum.Int).SetBig(p.Z1, p.Z1.BitLen())

	{
		// lhs = z₁ ⊙ C + rand
		lhs := public.C.Clone().Mul(verifier, z1)
		lhs.Randomize(verifier, new(safenum.Nat).SetBig(p.W, p.W.BitLen()))

		// rhsCt = A ⊕ (e ⊙ D)
		rhs := public.D.Clone().Mul(verifier, e).Add(verifier, p.A)

		if !lhs.Equal(rhs) {
			return false
		}
	}

	{
		// lhs = [z₁]G
		lhs := curve.NewIdentityPoint().ScalarBaseMult(curve.NewScalarBigInt(p.Z1))

		// rhs = [e]X + Bₓ
		rhs := curve.NewIdentityPoint().ScalarMult(curve.NewScalarInt(e), public.X)
		rhs.Add(rhs, &p.Bx)
		if !lhs.Equal(rhs) {
			return false
		}
	}

	return true
}

func challenge(hash *hash.Hash, public Public, commitment *Commitment) *safenum.Int {
	_ = hash.WriteAny(public.Aux, public.Verifier,
		public.C, public.D, public.X,
		commitment.A, commitment.Bx,
		commitment.E, commitment.S)

	return sample.IntervalScalar(hash.Digest())
}
