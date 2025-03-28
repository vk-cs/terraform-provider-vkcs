package main

//go:generate tfplugingen-framework generate resources --input vkcs/cdn/spec.json --output vkcs/cdn
//go:generate tfplugingen-framework generate data-sources --input vkcs/cdn/spec.json --output vkcs/cdn
