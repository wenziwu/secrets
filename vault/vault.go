package vault

import (
  "fmt"
  "strings"
  "os"
  "io/ioutil"
  "bufio"
  "encoding/json"
  "github.com/jarmo/secrets/secret"
)

func List(filter string) []secret.Secret {
  var secrets []secret.Secret
  for _, secret := range read() {
    if secret.Id.String() == filter ||
         strings.Index(strings.ToLower(secret.Name), strings.ToLower(filter)) != -1 {
      secrets = append(secrets, secret)
    }
  }

  return secrets
}

func Add(name string) secret.Secret {
  secrets := read()

  fmt.Println("Enter value:")
  var value []string
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
      value = append(value, scanner.Text())
  }

  secret := secret.Create(name, strings.Join(value, "\n"))
  write(append(secrets, secret))
  return secret
}

func read() []secret.Secret {
  if data, err := ioutil.ReadFile(storagePath()); os.IsNotExist(err) {
    return make([]secret.Secret, 0)
  } else {
    var secrets []secret.Secret
    if err := json.Unmarshal(data, &secrets); err != nil {
      panic(err)
    } else {
      return secrets
    }
  }
}

func write(secrets []secret.Secret) {
  secretsJSON, _ := json.Marshal(secrets)
  if err := ioutil.WriteFile(storagePath(), secretsJSON, 0600); err != nil {
    panic(err)
  }
}

func storagePath() string {
  return "/Users/jarmo/.secrets.json"
}
