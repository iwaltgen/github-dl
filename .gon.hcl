# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./dist/ghdl-macos_darwin_amd64/github-dl"]
bundle_id = "dev.iwaltgen.ghdl"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Apple Development: iwaltgen@gmail.com"
}
