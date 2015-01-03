# Description

Go implementation of yubikey server to be able to run your own server on network who can't have access to the official servers.

Store all information inside a sqlite database, need to create other connectors for different backend

# Usage

  $go build -> to build the server
  $./yubikey-server -app "NameOfYourApp" -> will add a new application and display the id and key
  $./yubikey-server -name YourName -pub "publicKey" -secret "AESSecret" -> will add a new key in the system"
  $./yubikey-server -s -> will start the server on the default port 3000