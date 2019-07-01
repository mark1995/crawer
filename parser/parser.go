package parser

import "errors"

type Parser interface {
	ParserUrl(url string) error
}

func NewParser(from string) (Parser, error) {
	switch from {
	case "booktxt":
		return new(BookTextSpider), nil
	default:
		return nil, errors.New("系统暂未处理该类型的配置文件")
	}
}
