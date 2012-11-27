package com.bypass;

import com.bodymedia.common.applets.CommandException;
import com.bodymedia.common.applets.device.util.LibraryException;
import com.bodymedia.device.serial.SerialPort3;
import com.bodymedia.device.usb.Usb;
import com.bodymedia.common.applets.logger.Logger;

import java.io.FileWriter;
import java.io.IOException;
import java.util.regex.Pattern;
import java.util.regex.Matcher;

public class BodyBuggBypass {
	public static void main(String[] args) throws CommandException, LibraryException, IOException {
		Logger log = Logger.getLogger();
		log.setPriority(3);
		
		Usb usb = new Usb("bmusbapex5", Usb.ANY);
		
		String[] ports = usb.getArmbandPorts();
		
		if(ports.length < 1) {
			System.out.println("No BodyBuggs detected.");
			System.exit(1);
		} else if(ports.length > 1) {
			System.out.println("Multiple BodyBuggs detected, re-run with only one connected.");
			System.exit(1);
		}
		
		String serialPort = ports[0];
		
		SerialPort3 ser = new SerialPort3("bmcommapex5", serialPort);
		ser.setAddr(0x0000000E, 0xFFFFFFFF);
		
		ser.open();
		
		ser.writeCommand("get lastdataupdate");
		
		Pattern lastUpdateRegex = Pattern.compile("Last Data Update: ([0-9]+)");
		Matcher matcher = lastUpdateRegex.matcher(ser.readResponse().toString());
		
		matcher.find();
		String lastUpdate = matcher.group(1);
		
		String logPath = String.format("%s.log", lastUpdate);
		System.out.printf("Writing data to: %s\n", logPath);
		FileWriter out = new FileWriter(logPath);
		
		ser.writeCommand("retrieve PDP");
		out.write(ser.readResponse().toString());
		
		out.close();
		
		System.out.println("Clearing device memory and updating timestamps.");
		ser.writeCommand("file init");
		ser.writeCommand(String.format("set lastdataupdate %d", System.currentTimeMillis() / 1000L));
		ser.writeCommand(String.format("set epoch %d", System.currentTimeMillis() / 1000L));
		
		ser.close();
	}
}
