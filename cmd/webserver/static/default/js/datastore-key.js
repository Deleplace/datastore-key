(function($) {

	//
	// Button handlers issue requests to services, retrieve JSON data and fill the form fields.
	//
	// Feel free to call /encode and /decode directly for your needs.
	//

	$("#ajax-encode").click(function(){
		$("#ajax-encode").button('loading');
		$.ajax({
			url: "/encode",
		    dataType: "json",
			data: $(".form-encode").serialize(),
			success: function(response) {
                let appID = $(".form-encode input[name=appid]").val();

                let oldbox = $(".form-decode-old textarea[name=oldkeystring]");
                if(appID && response.oldKeyString) {
					oldbox.val( response.oldKeyString );
					oldbox.focus();
					$("#ajax-encode").button('reset');
					$(".form-decode-old").effect("highlight", 1200);
                    $(".form-decode-old fieldset").effect("highlight", 1200);
				    window.history.pushState('', '', '/?oldkeystring=' + response.oldKeyString);
                } else {
					oldbox.val("");
                }

                let newbox = $(".form-decode-new textarea[name=newkeystring]");
                if(response.newKeyString) {
					newbox.val( response.newKeyString );
					newbox.focus();
					$("#ajax-encode").button('reset');
					$(".form-decode-new").effect("highlight", 1200);
                    $(".form-decode-new fieldset").effect("highlight", 1200);
                    if(!appID)
				        window.history.pushState('', '', '/?newkeystring=' + response.newKeyString);
                }else{
					newbox.val("");
                }
			},
			error: function(msg) {
			      alert( "Encoding went wrong : [" + err.responseText + "]" );
				  $("#ajax-encode").button('reset');
			}
		});
	});

	$.fn.ajaxDecode = function(flavor, afterSuccess){
        // flavor may be "old" of "new"
		$.ajax({
			url: "/decode",
		    dataType: "json",
			data: $(".form-decode-" + flavor).serialize(),
			success: function(response) {
				$(".form-encode").find("input[type=text], textarea").val("");  // Clear all values

				$(".form-encode input[name=kind]").val( response.kind );
				$(".form-encode input[name=intid]").val( response.intID );
				$(".form-encode input[name=stringid]").val( response.stringID );
				$(".form-encode input[name=appid]").val( response.appID );
				$(".form-encode input[name=appid]").change(); // update link
				$(".form-encode input[name=namespace]").val( response.namespace );

				if( response.parent ){
					setSectionVisibility( $("#set-parent"), $(".key-parent"), true, "Remove parent" );
					$(".form-encode input[name=kind2]").val( response.parent.kind );
					$(".form-encode input[name=intid2]").val( response.parent.intID );
					$(".form-encode input[name=stringid2]").val( response.parent.stringID );

					if( response.parent.parent ){
						setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), true, "Remove grandparent" );
						$(".form-encode input[name=kind3]").val( response.parent.parent.kind );
						$(".form-encode input[name=intid3]").val( response.parent.parent.intID );
						$(".form-encode input[name=stringid3]").val( response.parent.parent.stringID );
					}else{
						setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), false, "Set grandparent" );
					}
				}else{
					setSectionVisibility( $("#set-parent"), $(".key-parent"), false, "Set parent" );
					setSectionVisibility( $("#set-grand-parent"), $(".key-grand-parent"), false, "Set grandparent" );
				}
				if( afterSuccess )
					afterSuccess(response);
			},
			error: function(err) {
			      alert( "Key string seems invalid : [" + err.responseText + "]" );
			      $("#ajax-decode-" + flavor).button('reset');
			}
		});
	}

	$("#ajax-decode-old").click(function(){
        if( !$(".form-decode-old textarea[name=oldkeystring]").val()  )
            return;
		$("#ajax-decode-old").button('loading');
		$.fn.ajaxDecode( "old", function(response){
		    $("#ajax-decode-old").button('reset');
			$(".form-encode").effect("highlight", 1200);
			$(".form-encode fieldset").effect("highlight", 1200);
            window.history.pushState('', '', '/?oldkeystring=' + $(".form-decode-old textarea[name=oldkeystring]").val());

            // Decoding an old key will fill new key as well
            if(response.newkeystring){
                $(".form-decode-new textarea[name=newkeystring]").val(response.newkeystring);
                $(".form-decode-new").effect("highlight", 1200);
                $(".form-decode-new fieldset").effect("highlight", 1200);
            }
		});
    });

	$("#ajax-decode-new").click(function(response){
        if( !$(".form-decode-new textarea[name=newkeystring]").val()  )
            return;
		$("#ajax-decode-new").button('loading');
		$.fn.ajaxDecode( "new", function(){
		    $("#ajax-decode-new").button('reset');
			$(".form-encode").effect("highlight", 1200);
			$(".form-encode fieldset").effect("highlight", 1200);
            window.history.pushState('', '', '/?newkeystring=' + $(".form-decode-new textarea[name=newkeystring]").val());
            // Decoding a new key implies clearing the old key (which would lack AppID)
            $(".form-decode-old textarea[name=oldkeystring]").val("");
		});
	});

	//
	// Toggle Set/Remove the direct parent of the main key
	//
	$("#set-parent").click(function(){
		if( $(".key-parent").is(":visible") ){
			$(".form-encode .key-parent").find("input[type=text], textarea").val("");
			setSectionVisibility( $(this), $(".key-parent"), false, "Set parent" );
			setSectionVisibility( $(this), $(".key-grand-parent"), false, "Set grandparent" );
		}else{
			setSectionVisibility( $(this), $(".key-parent"), true, "Remove parent" );
		}
	});

	//
	// Toggle Set/Remove the parent of the parent of the main key
	//
	$("#set-grand-parent").click(function(){
		if( $(".key-grand-parent").is(":visible") ){
			$(".form-encode .key-grand-parent").find("input[type=text], textarea").val("");
			setSectionVisibility( $(this), $(".key-grand-parent"), false, "Set grandparent" );
		}else{
			setSectionVisibility( $(this), $(".key-grand-parent"), true, "Remove grandparent" );
		}
	});

	function setSectionVisibility(button, section, newVisibility, newButtonText){
		if( newVisibility ){
			section.removeClass("hidden");
			section.focus();
		}else{
			section.addClass("hidden");
		}
		button.html(newButtonText);
	}


	//
	// There is no support for great-grand-parents and further ancestors, but contact me if you feel you need that.
	//

	$(".form-encode input[name=appid]").change(function(){
		var appid = $(this).val();
		var url = "javascript:void(0);";
		if( appid ){
			var pos = appid.indexOf("~");
			if( pos != -1 )
				appid = appid.substring(pos+1);
			url = "https://" + appid + ".appspot.com";
		}
		$("#appspot-link").attr("href", url);
	});
	$(".form-encode input[name=appid]").change();

	$("#btn-more").click(function() {
	    $("#more-content").collapse('toggle');
	});

	$.fn.openInDatastoreViewer = function(){
	    var key = $(".form-decode-old textarea[name=oldkeystring]").val();
	    if( key ){
	    	var kind = $(".form-encode input[name=kind]").val();
	    	var appid = $(".form-encode input[name=appid]").val();
	    	var namespace = $(".form-encode input[name=namespace]").val();
	    	if( appid && kind ){
	    		var url= "https://appengine.google.com/datastore/explorer?submitted=1&app_id=" + appid
	    			+ "&show_options=yes&viewby=gql&query=SELECT+*+FROM+" + kind
	    			+ "+WHERE+__key__%3DKEY%28%27"+ key + "%27%29"
	    			+ "&namespace=" + namespace
	    			+ "&options=Run+Query" ;
	    		window.open( url, "datastoreViewer" );
	    	}else{
	    		alert("Please click the Decode button first, to retrieve the App ID.")
	    	}
	    }else{
    		alert("Please provide a datastore key");
    	}
	}

	$("#open-in-datastore-viewer").click($.fn.openInDatastoreViewer);

	$("#link-for-bookmark").click(function() {
	    var url= window.location.protocol + "//" + window.location.hostname + "/?";
	    [ "kind", "intid", "stringid", "appid", "namespace", "kind2", "intid2", "stringid2", "kind3", "intid3", "stringid3" ].forEach(function(f){
	    	var v = $(".form-encode input[name="+f+"]").val();
	    	if( v )
	    		url += f + "=" + encodeURIComponent(v) + "&";
	    });
	    var oldKeyString =  $(".form-decode-old textarea[name=oldkeystring]").val();
	    if( oldKeyString )
	    	url += "oldkeystring=" + encodeURIComponent(oldKeyString);
	    window.location = url;
	});

	$("#link-engine-this").click(function() {
		window.external.AddSearchProvider( "http://datastore-key.appspot.com/static/xml/opensearch-this.xml" );
	});

	$("#link-engine-ds-viewer").click(function() {
		window.external.AddSearchProvider( "http://datastore-key.appspot.com/static/xml/opensearch-jump-to-datastore-viewer.xml" );
	});


	$("#btn-about").click(function() {
		$("#about-content").collapse('toggle');
	});
})(jQuery);
