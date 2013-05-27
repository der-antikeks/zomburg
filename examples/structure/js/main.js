(function() {

	var SCREEN_WIDTH = window.innerWidth;
	var SCREEN_HEIGHT = window.innerHeight;

	var ws, wsAddr = window.wsaddr;
	var container, hud, stats;
	var camera, controls;
	var scene, renderer, projector, raycaster;
	var light;

	init();

	function init() {
		console.log("init");

		// check for webgl and websockets

		hud = document.getElementById("overlay");
		container = document.getElementById("main");

		// connect ws
		if (ws != null) {
			ws.close();
			ws = null;
		}
		
		console.log("connecting to", wsAddr);
		ws = new WebSocket("ws://"+wsAddr+"/ws");
		ws.onopen = onOpen;
		ws.onclose = onClose;
		ws.onmessage = onMessage;
		
		window.onbeforeunload = function() {
			ws.onclose = function() {};
			ws.close();
		};

		// load objects

		// scene
		// camera
		// lights
		// world
		// renderer

		// events
		window.addEventListener('resize', onWindowResize, false);
		document.addEventListener('keydown', onKeyDown, false);
		document.addEventListener('mousemove', onDocumentMouseMove, false);
		document.onselectstart = function() {return false};

		// start loop
		update();
	}

	function onOpen(e) {
		console.log("open", e);
	}
	
	function onClose(e) {
		console.log("close", e);
		ws = null;
		//window.setTimeout(connect, 3000);
	}
	
	function onMessage(e) {
		console.log("message", e);
		//var m = JSON.parse(e.data);
	}

	function onKeyDown(e) {
		switch(event.keyCode) {
			case 79: // O
			case 87: // W
			case 83: // S
			case 68: // D
			case 32: // Space
			case 17: // Ctrl
			default:
				console.log("key: " + e.keyCode);
		}
	}

	function onDocumentMouseMove(e) {
		e.preventDefault();

		/*
		mouse.x = (e.clientX / SCREEN_WIDTH) * 2 - 1;
		mouse.y = - (e.clientY / SCREEN_HEIGHT) * 2 + 1;
		*/
	}

	function onWindowResize() {
		SCREEN_WIDTH = window.innerWidth;
		SCREEN_HEIGHT = window.innerHeight;

		/*
		camera.aspect = SCREEN_WIDTH / SCREEN_HEIGHT;
		camera.updateProjectionMatrix();

		renderer.setSize(SCREEN_WIDTH, SCREEN_HEIGHT);
		*/
	}

	function update() {
		requestAnimationFrame(update);

		// pre render updates
		// controls, positions

		render();

		// post render updates
		// stats.update();
	}

	function render() {
		// renderer.render(scene, camera);
	}

})();
