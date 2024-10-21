package main

import "log"

func main() {
	if i := Init(); !i {
		log.Println("Failed to initialize AMD SMI")
		return
	}

	log.Println("AMD SMI initialized")

	sockets, err := GetSocketHandles()
	if err != nil {
		log.Println("failed to get socket handles:", err)
		return
	}

	log.Println("Socket handles: ", len(sockets))

	for _, socket := range sockets {
		name, err := GetSocketName(socket, 256)
		if err != nil {
			log.Println("failed to get socket name:", err)
			return
		}

		log.Println("Socket name:", name)

		processors, err := GetProcessorHandles(socket)
		if err != nil {
			log.Println("failed to get processor handles:", err)
			return
		}

		log.Println("Processors: ", len(processors))

		for _, processor := range processors {
			log.Println("=====================================")
			pType, err := GetProcessorType(processor)
			if err != nil {
				log.Println("failed to get processor type:", err)
				return
			}

			log.Println("Processor type:", pType)

			boardInfo, err := GetGPUBoardInfo(processor)
			if err != nil {
				log.Println("failed to get board info:", err)
				return
			}

			log.Println("Model Number: ", boardInfo.ModelNumber)
			log.Println("Serial Number: ", boardInfo.ProductSerial)
			log.Println("FRU ID: ", boardInfo.FruID)
			log.Println("Product Name: ", boardInfo.ProductName)
			log.Println("Manufacturer Name: ", boardInfo.ManufacturerName)

			gpuID, err := GetGPUID(processor)
			if err != nil {
				log.Println("failed to get GPU ID:", err)
				return
			}
			log.Println("ID: ", gpuID)

			vram, err := GetGPUVRAM(processor)
			if err != nil {
				log.Println("failed to get GPU VRAM:", err)
				return
			}
			log.Println("Total VRAM: ", vram.Total)
			log.Println("Used VRAM: ", vram.Used)

			uuid, err := GetGPUUUID(processor)
			if err != nil {
				log.Println("failed to get GPU UUID:", err)
				return
			}

			log.Println("UUID: ", uuid)
		}
	}

	if ok := Shutdown(); !ok {
		log.Println("Failed to shutdown AMD SMI")
		return
	}

	log.Println("AMD SMI shutdown")
}
