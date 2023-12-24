{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell
{
	nativeBuildInputs = [
		pkgs.go
		pkgs.cobra-cli
	];
}
