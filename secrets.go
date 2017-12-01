package main

import (
  "os"
  "fmt"
  "github.com/jarmo/secrets/cli"
  "github.com/jarmo/secrets/cli/command"
  "github.com/jarmo/secrets/vault"
  "github.com/jarmo/secrets/vault/storage"
  "github.com/jarmo/secrets/vault/storage/path"
  "github.com/jarmo/secrets/secret"
  "github.com/jarmo/secrets/input"
)

const VERSION = "0.0.1"

func main() {
  switch parsedCommand := cli.Execute(VERSION, os.Args[1:]).(type) {
    case command.List:
      secrets, _, _ := loadVault()
      for _, secret := range vault.List(secrets, parsedCommand.Filter) {
        fmt.Println(secret)
      }
    case command.Add:
      secrets, path, password := loadVault()
      secretName := parsedCommand.Name
      secretValue := input.AskMultiline(fmt.Sprintf("Enter value for '%s':\n", parsedCommand.Name))
      newSecret, newSecrets := vault.Add(secrets, secretName, secretValue)
      storage.Write(password, path, newSecrets)
      fmt.Println("Added:", newSecret)
    case command.Delete:
      secrets, path, password := loadVault()
      deletedSecret, newSecrets, err := vault.Delete(secrets, parsedCommand.Id)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      } else {
        storage.Write(password, path, newSecrets)
        fmt.Println("Deleted:", deletedSecret)
      }
    case command.Edit:
      secrets, path, password := loadVault()
      newName := input.Ask(fmt.Sprintf("Enter new name: "))
      newValue := input.AskMultiline("Enter new value:\n")
      editedSecret, newSecrets, err := vault.Edit(secrets, parsedCommand.Id, newName, newValue)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      } else {
        storage.Write(password, path, newSecrets)
        fmt.Println("Edited:", editedSecret)
      }
    case command.ChangePassword:
      currentPassword := askPassword()
      newPassword := input.AskPassword("Enter new vault password: ")
      newPasswordConfirmation := input.AskPassword("Enter new vault password again: ")

      if err := vault.ChangePassword(path.Get(), currentPassword, newPassword, newPasswordConfirmation); err != nil {
        fmt.Println(err)
        os.Exit(1)
      } else {
        fmt.Println("Vault password successfully changed!")
      }
    default:
      fmt.Printf("Unhandled command: %T\n", parsedCommand)
  }
}

func loadVault() ([]secret.Secret, string, []byte) {
  password := askPassword()
  vaultPath := path.Get()
  return storage.Read(password, vaultPath), vaultPath, password
}

func askPassword() []byte {
  return input.AskPassword("Enter vault password: ")
}
