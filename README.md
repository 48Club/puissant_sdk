# puissant_sdk

## Usage SDK

### get package

```bash
> go get -u github.com/bnb48club/puissant_sdk/bnb48.sdk
```

### import package

```go
import "github.com/bnb48club/puissant_sdk/bnb48.sdk"
```

### example

```go
package main

import (
	"context"
	"log"

	"github.com/bnb48club/puissant_sdk/bnb48.sdk"
)

func main() {
	client, err := bnb48.Dial("general rpc host", "puissant rpc host")
	if err != nil {
		log.Panicln(err.Error())
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Printf("gasPrice: %s", gasPrice.String())
}
```
## demo require

- bsc mainnet account balance is greater than 0.1 tbnb

### test demo selfTransfer

#### create `config.yaml` file

```yaml
wallet:
  # wallet privateKey
  - d88873b743af4da75edc17c0430be1691e1a0801b2a76575bc85597b6c17e336
```

#### send puissant txs
```bash
> git clone https://github.com/bnb48club/puissant_sdk.git && cd puissant_sdk
> go run github.com/bnb48club/puissant_sdk/demo/selfTransfer
```

### test demo multiTransfer

#### create `config.yaml` file

```yaml
wallet:
  # wallet privateKeys
  - d88873b743af4da75edc17c0430be1691e1a0801b2a76575bc85597b6c17e336
  - 17cddb8204ed4e808aef43a97a6f20dec8070426ce967c7375bdd33fd0693807
  - be33dfca66b35bc107e5d8e7f6e0f6232e81149b7501c36cbdbd9c52f571202c
  - 2d79317acd39955d97249c9b942031b645df61fdcadd328cb487387840c6b086
  - 7a82ed6033f5e24074021acbbb8a8dc90b87a6d67e4c18b1634f328dbc34cd29
```

#### send puissant txs
```bash
> git clone https://github.com/bnb48club/puissant_sdk.git && cd puissant_sdk
> go run github.com/bnb48club/puissant_sdk/demo/multiTransfer
```