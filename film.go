package main

type Films struct{
	films []Film `json:Films`
}

type Film struct{

	Title string `json:title`
	Director string `json:director`
}