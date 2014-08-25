package spaceballoon

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

var options = serial.OpenOptions{
	PortName:        "/dev/ttyUSB0",
	BaudRate:        9600,
	DataBits:        8,
	StopBits:        1,
	MinimumReadSize: 4,
}

type Dongle struct {
	port     io.ReadWriteCloser
	phoneNum string
	timeout  time.Duration
}

func NewDongle(phoneNum string) (*Dongle, error) {
	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	return &Dongle{port, phoneNum, time.Second}, err
}

func (d *Dongle) Stop() {
	d.port.Close()
}

func (d *Dongle) Test() error {
	err := d.write(smsTest)
	if err != nil {
		return err
	}

	return err
}

const maxSMSLength = 160

const (
	smsTest        = "AT\r\n"
	smsTextMode    = "AT+CMGF=1\r\n"
	smsSetPhoneNum = "AT+CMGS=\"%s\"\r\n"
	smsEndBody     = "\x1A"
)

func (d *Dongle) SendSMS(message []string) error {
	var totalLen int
	for _, line := range message {
		totalLen += len(line)
	}
	if totalLen > maxSMSLength {
		return errors.New("SMS too long")
	}

	// Setting SMS mode
	if err := d.write(smsTextMode); err != nil {
		return err
	}

	// Setting phone number
	if err := d.write(fmt.Sprintf(smsSetPhoneNum, d.phoneNum)); err != nil {
		return err
	}

	// Sending SMS
	for i, line := range message {
		// Writing line
		var err error
		if i < len(message)-1 {
			err = d.write(line + "\r\n")
		} else {
			err = d.write(line)
		}
		if err != nil {
			return err
		}
	}

	// Finishing
	if err := d.write(smsEndBody); err != nil {
		return err
	}
	return nil
}

func (d *Dongle) write(text string) error {
	_, err := d.port.Write([]byte(text))
	if err != nil {
		return err
	}

	time.Sleep(d.timeout)

	bytes := make([]byte, len(text)+60)
	_, err = d.port.Read(bytes)
	return err
}
