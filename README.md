# ddc

A Linux implementation for communicating with a display using the dcc protocol over i2c for getting and setting values like brightness or contrast. 

I have used ddcutil as inspiration and the source for better understanding the protocols. But for my project which was BluetoothBLE server with GATT that does a passthrough on the read and write handler this library was bit to much overhead and didn't provide me with the things i needed.

This is a low level implementation so don`t expect any support for probing, feature translations or tables and is only made to quickly write/read features or get edid info from connected screens.

To get information about your monitor capabilities i suggest to have a look at [`ddcutil`](https://www.ddcutil.com/).

https://milek7.pl/ddcbacklight/mccs.pdf