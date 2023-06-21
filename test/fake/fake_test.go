package fake_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"x-gwi/app/client"
	"x-gwi/app/logs"
	userAPIv2 "x-gwi/proto/serv/_user/v2"
	assetAPIv2 "x-gwi/proto/serv/asset/v2"
	favouriteAPIv2 "x-gwi/proto/serv/favourite/v2"
	opinionAPIv2 "x-gwi/proto/serv/opinion/v2"
	"x-gwi/service"
	"x-gwi/test/fake"
)

func Example_protojson_Marshal() {
	var i int64 = 25
	t := time.Now()
	u := fake.FakeUserCore(i)
	b, _ := protojson.Marshal(u)
	fmt.Println(service.NameUser)
	fmt.Println(string(b))

	a := fake.FakeAssetCore(1)
	b, _ = protojson.Marshal(a)
	fmt.Println(service.NameAsset)
	fmt.Println(string(b))

	o := fake.FakeFavouriteCore(1, 1, 1)
	b, _ = protojson.Marshal(o)
	fmt.Println(service.NameFavourite)
	fmt.Println(string(b))

	f := fake.FakeOpinionCore(1, 1, 1)
	b, _ = protojson.Marshal(f)
	fmt.Println(service.NameOpinion)
	fmt.Println(string(b))

	_ = b
	fmt.Println(time.Since(t))

	// Output:
}

func Example_json_Marshal() {
	var i int64 = 1446
	t := time.Now()
	u := fake.FakeUserCore(i)
	b, _ := json.Marshal(u)
	fmt.Println(service.NameUser)
	fmt.Println(string(b))

	a := fake.FakeAssetCore(1)
	b, _ = json.Marshal(a)
	fmt.Println(service.NameAsset)
	fmt.Println(string(b))

	o := fake.FakeFavouriteCore(1, 1, 1)
	b, _ = json.Marshal(o)
	fmt.Println(service.NameFavourite)
	fmt.Println(string(b))

	f := fake.FakeOpinionCore(1, 1, 1)
	b, _ = json.Marshal(f)
	fmt.Println(service.NameOpinion)
	fmt.Println(string(b))

	_ = b
	fmt.Println(time.Since(t))

	// Output:
}

func Example_gRPC_Client_loading_fake_data() {
	const (
		timeoutgRPC = 5 * time.Minute
	)

	grpcClientConn, cancelDial, err := client.InsecureClientConnGRPC(context.Background(), client.NewConfigClient())
	defer cancelDial()
	if err != nil {
		logs.Error().Err(err).Send()
		return
	}

	// userAPIv2 "x-gwi/proto/serv/_user/v2"
	// assetAPIv2 "x-gwi/proto/serv/asset/v2"
	// assetAPIv2 "x-gwi/proto/serv/asset/v2"
	// favouriteAPIv2 "x-gwi/proto/serv/favourite/v2"
	// opinionAPIv2 "x-gwi/proto/serv/opinion/v2"
	cliU2 := userAPIv2.NewUserServiceClient(grpcClientConn)
	cliA2 := assetAPIv2.NewAssetServiceClient(grpcClientConn)
	cliF2 := favouriteAPIv2.NewFavouriteServiceClient(grpcClientConn)
	cliO2 := opinionAPIv2.NewOpinionServiceClient(grpcClientConn)

	ctxRPC, cancelRPC := context.WithTimeout(context.Background(), timeoutgRPC)
	defer cancelRPC()

	const (
		xu int64 = 1_000
		xa int64 = 4_000
		xf int64 = 2_000
		xo int64 = 2_000
	)

	var (
		iu, ia, eu, ea, nu, na int64
	)

	wg := new(sync.WaitGroup)

	t := time.Now()

	go func() {
		wg.Add(1)
		// u:
		for iu = 1; iu <= xu; iu++ {
			_, err := cliU2.Create(ctxRPC, fake.FakeUserCore(iu))
			if err != nil {
				// logs.Error().Err(err).Int64("iu", iu).Send()
				eu++
				// if eu > 5 {
				// 	break u
				// }
			} else {
				nu++
				// b, _ := protojson.Marshal(out)
				// fmt.Println(service.NameUser)
				// fmt.Println(string(b))
			}

		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		// a:
		for ia = 1; ia <= xa; ia++ {
			_, err := cliA2.Create(ctxRPC, fake.FakeAssetCore(ia))
			if err != nil {
				ea++
			} else {
				na++
			}
		}
		wg.Done()
	}()

	runtime.Gosched()
	wg.Wait()

	fmt.Printf("iu: %d, u: %d, eu: %d, ia: %d, a: %d, ea: %d, t: %v\n", iu-1, nu, eu, ia-1, na, ea, time.Since(t))

	var (
		ifa, iop, efa, eop, nfa, nop int64
	)

	t2 := time.Now()

	go func() {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		// fa:
		for ifa = 1; ifa <= xf; ifa++ {
			ru := rand.Int63n(xu) + 1
			ra := rand.Int63n(xa) + 1
			_, err := cliF2.Create(ctxRPC, fake.FakeFavouriteCore(ifa, ru, ra))
			if err != nil {
				// "favouriteCore.Create: c.edgeKeyFromTo: favourite key already exists"
				efa++
			} else {
				nfa++
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		// op:
		for iop = 1; iop <= xo; iop++ {
			ru := rand.Int63n(xu) + 1
			ra := rand.Int63n(xa) + 1
			_, err := cliO2.Create(ctxRPC, fake.FakeOpinionCore(iop, ru, ra))
			if err != nil {
				eop++
			} else {
				nop++
			}
		}
		wg.Done()
	}()

	runtime.Gosched()
	wg.Wait()

	fmt.Printf("ifa: %d, fa: %d, efa: %d, iop: %d, op: %d, eop: %d, t2: %v, t: %v\n", ifa-1, nfa, efa, iop-1, nop, eop, time.Since(t2), time.Since(t))

	// Output:
}
