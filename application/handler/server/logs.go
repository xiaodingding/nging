/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package server

import (
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/tail"
	"github.com/webx-top/echo"
)

func LogShow(ctx echo.Context) error {
	var logFile string
	typ := ctx.Paramx(`type`).String()
	switch typ {
	case `caddy`:
		logFile = config.DefaultConfig.Caddy.LogFile
	}
	data := ctx.Data()
	if len(logFile) == 0 {
		data.SetData(`没有日志文件`)
		return ctx.JSON(data)
	}
	lastLines := ctx.Formx(`lastLines`).Int()
	config := tail.Config{
		MustExist: true,
		LastLines: 50,
	}
	if lastLines > 0 {
		config.LastLines = lastLines
	}
	obj, err := tail.TailFile(logFile, config)
	if err != nil {
		data.SetError(err)
	} else {
		var content string
		for line := range obj.Lines {
			content += line.Text + "\n"
		}
		data.SetData(content)
	}
	return ctx.JSON(data)
}
