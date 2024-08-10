window.addEventListener("load", function load(event){
    
    // Null Island
    var map = L.map('map').setView([0.0, 0.0], 12);
    
    var tile_layer = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
	maxZoom: 19,
    });
    
    tile_layer.addTo(map);

    fetch("/features.geojson")
	.then((rsp) => rsp.json())
	.then((f) => {

	    var raw_el = document.querySelector("#raw");

	    var format = function(str){

		// Remember: wof_format is defined by the /wasm/wof_format.wasm binary.
		// Details below.
		    
		wof_format(str).then((rsp) => {
		    append(rsp);
		}).catch((err) => {
		    console.log("Unable to format feature", err, str);
		    append(str);
		});
	    };
	    
	    var append = function(str) {
		var pre = document.createElement("pre");
		pre.appendChild(document.createTextNode(str));		    
		raw_el.appendChild(pre);
	    };

	    /*
	    if (raw_el){

		// Remember: Both sfomuseum.wasm.fetch and the WASM binary are imported and registered
		// in show.go. For details see: https://github.com/whosonfirst/go-whosonfirst-format-wasm
		
		sfomuseum.wasm.fetch("/wasm/wof_format.wasm").then(rsp => {

		    var features = f.features;
		    var count = features.length;

		    for (var i=0; i < count; i++){
			var str_f = JSON.stringify(features[i], "", " ");		    			
			format(str_f);
		    }
		    
		}).catch((err) => {
		    console.log("Unable to load wof_format.wasm", err);
		    var str_f = JSON.stringify(f, "", " ");		    
		    append(str_r);
		});

	    }
	    */
	    
	    var pt_handler = whosonfirst.spelunker.leaflet.handlers.point({});

	    var poly_style = whosonfirst.spelunker.leaflet.styles.consensus_polygon();	    
	    // var lbl_style = whosonfirst.spelunker.leaflet.styles.label_centroid();
	    
	    var geojson_args = {
		style: poly_style,
		pointToLayer: pt_handler,		
	    };

	    var geojson_layer = L.geoJSON(f);
	    geojson_layer.addTo(map);

	    var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);

	    var sw = bounds[0];
	    var ne = bounds[1];

	    if ((sw[0] == ne[0]) && (sw[1] == ne[1])){
		map.setView(sw, 12);
	    } else {
		map.fitBounds(bounds);
	    }
	    
	}).catch((err) => {
	    console.log("SAD", err);
	});
    
});