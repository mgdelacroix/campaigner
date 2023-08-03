{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    gnumake
    go_1_18
    golangci-lint
  ];
  shellHook = ''
    echo '
┏━╸┏━┓┏┳┓┏━┓┏━┓╻┏━╸┏┓╻┏━╸┏━┓
┃  ┣━┫┃┃┃┣━┛┣━┫┃┃╺┓┃┗┫┣╸ ┣┳┛
┗━╸╹ ╹╹ ╹╹  ╹ ╹╹┗━┛╹ ╹┗━╸╹┗╸
'
  '';
  hardeningDisable = [ "all" ];
}
