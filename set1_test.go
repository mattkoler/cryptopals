package main

import "testing"

/*
func TestSet1(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "s1c1",
			input:  "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d",
			output: "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
*/

func TestSet1C1(t *testing.T) {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	want := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	got, err := hexTo64(input)
	if err != nil {
		t.Fatalf("hex to string failed: %s", err)
	}
	if got != want {
		t.Errorf("hexTo64(%q) = %q, want %q", input, got, want)
	}
}

func TestSet1C2(t *testing.T) {
	input := "1c0111001f010100061a024b53535009181c"
	input2 := "686974207468652062756c6c277320657965"
	want := "746865206b696420646f6e277420706c6179"
	got, err := fixedXOR(input, input2)
	if err != nil {
		t.Fatalf("fixed xor failed: %s", err)
	}
	if got != want {
		t.Errorf("fixedXOR(%q, %q) = %q, want %q", input, input2, got, want)
	}
}

func TestSet1C3(t *testing.T) {
	input := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	want := "Cooking MC's like a pound of bacon"
	got, _, _, err := singleByteXOR(input)
	if err != nil {
		t.Fatalf("single byte xor failed: %s", err)
	}
	if got != want {
		t.Errorf("singleByteXOR(%q) = %q, want %q", input, got, want)
	}
}

func TestSet1C4(t *testing.T) {
	input := "./set1c4.txt"
	want := "Now that the party is jumping\n"
	got, _, err := findXORstring(input)
	if err != nil {
		t.Fatalf("find xor failed: %s", err)
	}
	if got != want {
		t.Errorf("singleByteXOR(%q) = %q, want %q", input, got, want)
	}
}

func TestSet1C5(t *testing.T) {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"
	want := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	got, err := repeatingKeyXOR(input, key)
	if err != nil {
		t.Fatalf("reapting xor failed: %s", err)
	}
	if got != want {
		t.Errorf("repeatingKeyXOR(%q, %q) = %q, want %q", input, key, got, want)
	}
}

func TestHamming(t *testing.T) {
	input := "this is a test"
	input2 := "wokka wokka!!!"
	want := 37
	got, err := hammingDist(input, input2)
	if err != nil {
		t.Fatalf("hamming dist failed: %s", err)
	}
	if got != want {
		t.Errorf("hammingDist(%q, %q) = %q, want %q", input, input2, got, want)
	}
}

func TestFindKeySize(t *testing.T) {
	input := "./sec1c6.txt"
	want := []int{2, 5, 29}
	got, err := findKeysize(input)
	if err != nil {
		t.Fatalf("find key size failed: %s", err)
	}
	for i := 0; i < len(want); i++ {
		if got[i] != want[i] {
			t.Errorf("findKeySize(%q) = %d, want %d", input, got, want)
		}
	}
}

func TestTransposeBlocks(t *testing.T) {
	input := "./sec1c6.txt"
	size := 5
	got, err := transposeBlocks(input, size)
	if err != nil {
		t.Fatalf("transpose blocks failed: %s", err)
	}
	if len(got) != size {
		t.Errorf("transposeBlocks(%q, %d) = %d, want %d", input, size, len(got), size)
	}
}

func TestSolveRepeatingKey(t *testing.T) {
	input := "./sec1c6.txt"
	want := ""
	got, err := solveRepeatingKey(input)
	if err != nil {
		t.Fatalf("solve repeating failed: %s", err)
	}
	if got != want {
		t.Errorf("solveRepeatingKey(%s) = %s, want %s", input, got, want)
	}
}
