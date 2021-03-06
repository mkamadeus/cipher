package hill

import (
	"errors"
	"fmt"
	"math"

	"github.com/mkamadeus/cipher/common/stringutils"
	"gonum.org/v1/gonum/mat"
)

func Decrypt(cipher string, key string) (string, error) {
	cipher = stringutils.Normalize(cipher)
	key = stringutils.Normalize(key)

	if !isQuadratic(len(key)) {
		return "", errors.New("key length isn't quadratic")
	}
	//Key Matrix
	keyMatrix := BuildKeyMatrix(key)
	keyDim := len(keyMatrix[0])
	keyDenseMatrix := mat.NewDense(keyDim, keyDim, MatToFloatArr(keyMatrix))

	// Determinan
	detKeyDense := mat.Det(keyDenseMatrix)
	detKeyRounded := int(math.Round(detKeyDense))
	fmt.Println(detKeyRounded)
	inversDetMod := ModInverse(detKeyRounded, 26)
	if inversDetMod == -1 {
		return "", errors.New("mod 26 key determinant not found")
	}
	// Invers
	var keyInv mat.Dense
	err := keyInv.Inverse(keyDenseMatrix)
	if err != nil {
		return "", errors.New("Key isn't invertible")
	}
	// I * D so it will fit with the formula
	var normalizedKeyInv mat.Dense
	normalizedKeyInv.Scale(detKeyDense, &keyInv)

	// Rounding
	a, b := normalizedKeyInv.Dims()
	normalizedRoundedKeyInv := make([][]int, a)
	for k := 0; k < a; k++ {
		normalizedRoundedKeyInv[k] = make([]int, b)
		for l := 0; l < b; l++ {
			normalizedRoundedKeyInv[k][l] = int(math.Round(normalizedKeyInv.At(k, l)))
		}
	}

	// Make new Dense
	normalizedKeyInversDense := mat.NewDense(a, b, MatToFloatArr(normalizedRoundedKeyInv))
	// Multiply with modInverseDet
	var finalKeyDense mat.Dense
	finalKeyDense.Scale(float64(inversDetMod), normalizedKeyInversDense)

	segmentedWord := BuildSegmentedPlainMat(BuildPlainString(cipher, keyDim))
	transposedSegment := transpose(segmentedWord)
	denseTransposedSegment := mat.NewDense(len(transposedSegment), len(transposedSegment[0]), MatToFloatArr(transposedSegment))

	var resMatrix mat.Dense
	resMatrix.Mul(&finalKeyDense, denseTransposedSegment)
	transposedRes := resMatrix.T()

	r, c := transposedRes.Dims()
	runeResMatrix := make([]rune, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			runeResMatrix[i*c+j] = rune(CorrectModulus(int(transposedRes.At(i, j)), 26)) + 65
		}
	}

	return string(runeResMatrix), nil
}
