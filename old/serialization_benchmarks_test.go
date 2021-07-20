package protobench

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	proto1jsonpb "github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString(l int) string {
	buf := make([]byte, l)
	for i := 0; i < (l+1)/2; i++ {
		buf[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", buf)[:l]
}

// github.com/golang/protobuf (aka V1old)

func generateGoV1(n int) []*GoV1 {
	a := make([]*GoV1, 0, n)
	for i := 0; i < n; i++ {
		a = append(a, &GoV1{
			Name:     randString(16),
			BirthDay: time.Now().UnixNano(),
			Phone:    randString(10),
			Siblings: rand.Int31n(5),
			Spouse:   rand.Intn(2) == 1,
			Money:    rand.Float64(),
			Type:     TypeV1(rand.Intn(4)),
			Values:   &GoV1_ValueS{ValueS: randString(5)},
			Address: &AddressV1{
				Street:   randString(5),
				Number:   rand.Int31n(20),
				PostCode: rand.Int31n(5),
				Floor:    rand.Int31n(5),
				Random: &RandomV1{
					FieldA: randString(5),
					FieldB: randString(5),
					FieldC: randString(5),
					FieldD: randString(5),
					FieldE: randString(5),
					FieldF: randString(5),
					FielG:  randString(5),
					FieldH: randString(5),
					FieldI: randString(5),
					FieldJ: randString(5),
					FieldK: rand.Int63n(5),
					FieldL: rand.Int63n(5),
					FieldM: rand.Int63n(5),
					FieldN: rand.Int63n(5),
					FieldO: randString(5),
					FieldP: randString(5),
					FieldQ: rand.Int31n(5),
					FieldR: randString(5),
					FieldS: randString(5),
					FieldT: randString(5),
					FieldU: rand.Int31n(5),
					FieldV: rand.Int31n(5),
					FieldW: rand.Int31n(5),
					FieldX: rand.Int31n(5),
					FieldY: randString(5),
					FieldZ: rand.Intn(2) == 1,
					NestedRandom: &NestedRandomV1{
						FieldA: randString(5),
						FieldB: randString(5),
						FieldC: randString(5),
						FieldD: randString(5),
						FieldE: randString(5),
						FieldF: randString(5),
						FielG:  rand.Float64(),
						FieldH: rand.Float64(),
						FieldI: rand.Float64(),
						FieldJ: rand.Float64(),
						FieldK: rand.Float64(),
					},
				},
			},
		})
	}
	return a
}

func Benchmark_GoV1old_Proto_Marshal(b *testing.B) {
	data := generateGoV1(b.N)
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		bytes, err := proto.Marshal(data[rand.Intn(len(data))])
		if err != nil {
			b.Fatal(err)
		}
		serialSize += len(bytes)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func Benchmark_GoV1old_Proto_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateGoV1(b.N)
	ser := make([][]byte, len(data))
	var serialSize int
	for i, d := range data {
		var err error
		ser[i], err = proto.Marshal(d)
		if err != nil {
			b.Fatal(err)
		}
		serialSize += len(ser[i])
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n := rand.Intn(len(ser))
		o := &GoV1{}
		err := proto.Unmarshal(ser[n], o)
		if err != nil {
			b.Fatalf("goprotobuf failed to unmarshal: %s (%s)", err, ser[n])
		}
	}
}
func Benchmark_GoV1old_JSON_Marshal(b *testing.B) {
	data := generateGoV1(b.N)
	marshaler := proto1jsonpb.Marshaler{}
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := marshaler.Marshal(&buf, data[rand.Intn(len(data))])
		if err != nil {
			b.Fatal(err)
		}
		serialSize += buf.Len()
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func Benchmark_GoV1old_JSON_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateGoV1(b.N)
	marshaler := proto1jsonpb.Marshaler{}
	ser := make([]bytes.Buffer, len(data))
	var serialSize int
	for i, d := range data {
		err := marshaler.Marshal(&ser[i], d)
		if err != nil {
			b.Fatal(err)
		}
		serialSize += ser[i].Len()
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	unmarshaler := proto1jsonpb.Unmarshaler{}
	randomI := randomI(b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n := randomI[i]
		o := &GoV1{}
		err := unmarshaler.Unmarshal(&ser[n], o)
		if err != nil {
			b.Fatalf("goprotobuf failed to unmarshal: %s (%s)", err, ser[n].String())
		}
	}
}

func randomI(n int) []int {
	randomI := make([]int, n)
	for i := 0; i < len(randomI); i++ {
		randomI[i] = i
	}
	rand.Shuffle(len(randomI), func(i, j int) {
		randomI[i], randomI[j] = randomI[j], randomI[i]
	})
	return randomI
}
