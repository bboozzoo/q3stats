$(document).ready(
    function() {
	$("[title]").tooltip(
	    {
		effect: "fade",
		offset: [-2, 0]
	    }
	);
	$("ul.tabs").tabs("div.panes > div");
    }
);