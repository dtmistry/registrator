# Registrator

[Registrator](https://github.com/gliderlabs/registrator) with an additional BigIp backend.

The BigIp backend uses the [iControl](https://devcentral.f5.com/login?returnurl=%2fwiki%2fiControlREST.HomePage.ashx) REST Api to add/remove services to a pool.

You can also provide an image name at startup and registrator will only act on docker events from containers of that image.

## Usage

Edit the DEV_RUN_OPTS in the 'make' file and simply run make. 

	DEV_RUN_OPTS ?= -service <image-name> bigip://<user>:<pass>@<bigip-host>/<bigip-pool>

Alternatively, run the below command

	$ docker run -d \
		--volume:/var/run/docker.sock:/tmp/docker.sock \
		--net=host \
		dmistry/registrator \
			-service <image-name> bigip://<user>:<pass>@<bigip-host>/<bigip-pool>


*image-name* = Name of the container image registrator should act on (Optional)  
*user* = BigIp user who has access to modify configuration objects  
*pass* = BigIp user password  
*bigip-host* = BigIp management host  
*bigip_pool* = BigIp pool where services will be added  


## License

MIT
