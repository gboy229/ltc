package team

import "hotgo/addons/users/service"

type sTeam struct{}

func NewTeam() *sTeam { return &sTeam{} }

func init() { service.RegisterTeam(NewTeam()) }
