package zk

import (
	"fmt"
	"math/big"

	"github.com/cronokirby/safenum"
	"github.com/taurusgroup/cmp-ecdsa/pkg/paillier"
	"github.com/taurusgroup/cmp-ecdsa/pkg/pedersen"
)

var (
	ProverPaillierPublic   *paillier.PublicKey
	ProverPaillierSecret   *paillier.SecretKey
	VerifierPaillierPublic *paillier.PublicKey
	VerifierPaillierSecret *paillier.SecretKey
	Pedersen               *pedersen.Parameters
)

func generate() {
	sk1 := paillier.NewSecretKey()
	sk2 := paillier.NewSecretKey()
	fmt.Printf("p1, _ := new(safenum.Nat).SetHex(\"%s\")\n", sk1.P().Hex()[2:])
	fmt.Printf("q1, _ := new(safenum.Nat).SetHex(\"%s\")\n", sk1.Q().Hex()[2:])
	fmt.Printf("p2, _ := new(safenum.Nat).SetHex(\"%s\")\n", sk2.P().Hex()[2:])
	fmt.Printf("q2, _ := new(safenum.Nat).SetHex(\"%s\")\n", sk2.Q().Hex()[2:])

	fmt.Println("ProverPaillierSecret = paillier.NewSecretKeyFromPrimes(p1, q1)")
	fmt.Println("VerifierPaillierSecret = paillier.NewSecretKeyFromPrimes(p2, q2)")
	fmt.Println("ProverPaillierPublic = ProverPaillierSecret.publicKey()")
	fmt.Println("VerifierPaillierPublic = VerifierPaillierSecret.publicKey()")
	ped, _ := sk2.GeneratePedersen()
	fmt.Printf("s, _ := new(big.Int).SetString(\"%s\", 10)\n", ped.S())
	fmt.Printf("t, _ := new(big.Int).SetString(\"%s\", 10)\n", ped.T())
	fmt.Println("Pedersen = &pedersen.Parameters{N: VerifierPaillierPublic.N, S: s, T: t}")
}

func init() {
	p1, _ := new(safenum.Nat).SetHex("FD90167F42443623D284EA828FB13E374CBF73E16CC6755422B97640AB7FC77FDAF452B4F3A2E8472614EEE11CC8EAF48783CE2B4876A3BB72E9ACF248E86DAA5CE4D5A88E77352BCBA30A998CD8B0AD2414D43222E3BA56D82523E2073730F817695B34A4A26128D5E030A7307D3D04456DC512EBB8B53FDBD1DFC07662099B")
	q1, _ := new(safenum.Nat).SetHex("DB531C32024A262A0DF9603E48C79E863F9539A82B8619480289EC38C3664CC63E3AC2C04888827559FFDBCB735A8D2F1D24BAF910643CE819452D95CAFFB686E6110057985E93605DE89E33B99C34140EF362117F975A5056BFF14A51C9CD16A4961BE1F02C081C7AD8B2A5450858023A157AFA3C3441E8E00941F8D33ED6B7")
	p2, _ := new(safenum.Nat).SetHex("EEEFE9909452DEC61592452661DA397DB0A0A2BCFD7F6FC07EFDF98DAAE3BA276AA244E3162E95196E87BD73902EDF9F3823C90E239E683E37973185B30746D06CD901581448F0FEF3EFEDDD5DD21904ED3EEB6C8381ABFF6A3F41CD0B1ADD61F9E74DCA871404AE2813FFAA1886FEE6F896F647D8F296877F8F728D008C4CB3")
	q2, _ := new(safenum.Nat).SetHex("F5D4C0FE31CE43EFC25BC3D8AA3E14C8F9D831A932ABEFB9755A7A0556BB9F6CA63C1CB242703FEA151952888C37850EF5D024BA2D6A0D081196FF1616DE9BB7E7FFD13D857CA9C382896E9B772C2C5358EC99A50505DF52F98BC7EB7215912F2CECBE723BF261A87F1F13D9964A3B318FD60AAE176D5A37A4855F5C0E6D7F83")
	ProverPaillierSecret = paillier.NewSecretKeyFromPrimes(p1, q1)
	VerifierPaillierSecret = paillier.NewSecretKeyFromPrimes(p2, q2)
	ProverPaillierPublic = ProverPaillierSecret.PublicKey
	VerifierPaillierPublic = VerifierPaillierSecret.PublicKey
	s, _ := new(big.Int).SetString("25448182756540222866319501898634270977599693736075116251767811398321761734336492575194265315792355808584718336663782672204688106417478368963300473751214591342337476683702712849651809231550640975602518619906861486655178732000388791254150785099820854279346548346528112423491316470201233748237839902669525854602668864974189022599379366745615480940437093919632552402353187455676393853230065384108683249537894167198489453490494046634358648195677341022149808031953931584669015124559681899432749131037714975312907114971572056736764605968515055562566021058441965744442163176054557494286204809259988469497547181346009783327359", 10)
	t, _ := new(big.Int).SetString("8015316201083753856167999987758855473596936298647823494052059806385270716210155885167716563507417164942757686623899459374045961056288024798233083861233712530311639442975277876900958870072552064612806429411065226076877525019560729216371707422331329043485629581836713694830988754061121605813376665575545793435480194071036877604266918397803979335525965949555983633053355236016000661904600828377166630829936488123935666148193019631612745852632950704726062001705644025539422488832473576716710166218032931524855280744450703912325492865050739662260518593977876485434650101403209615051036271024149130698937966662896144492264", 10)
	Pedersen, _ = pedersen.New(VerifierPaillierPublic.N(), s, t)
}
