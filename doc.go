/*

Travel is an HTTP router that provides routing similar to the "traversal" system from the Pyramid web framework in Python.

For details on the original traversal system please read: http://docs.pylonsproject.org/docs/pyramid/en/latest/narr/traversal.html

Simply put, traversal allows you to dynamically route HTTP requests by providing a nested map[string]interface{} object called the
"root tree" (Pyramid calls this the *resource tree*). Request URLs are tokenized and recursive lookup is performed on the root
tree object.

Example:

If the request URL is "/foo/bar/baz/123", it is tokenized to the following:

   ["foo", "bar", "baz", "123"]

The equivalent of the following lookup is performed:

   root_tree["foo"]["bar"]["baz"]["123"]

The object that results from this lookup is the "current object" (Pyramid calls this the "context"--travel provides a context object that
contains the current object).

How this maps to handlers (in Pyramid terminology: "views") depends upon the options passed when creating the router. Under traditional
traversal, if the lookup fully succeeded (no missing key errors), the name of the handler would be the empty string ("") which is considered
the default handler. If the lookup failed at any point, the handler name would be the token that failed and any remainder of the URL would
be passed to the handler as the "subpath" (see original traversal documentation linked above for more details).

Usage:

	func defaultHandler(w http.ResponseWriter, r *http.Request, c *travel.Context) {
		// handler code here
	}

	func errorHandler(w http.ResponseWriter, r *http.Request, err travel.TraversalError) {
		// HTTP error handler code here
		http.Error(w, err.Error(), err.Code())
	}

	func getRootTree() {
	  // Fetch root tree here
	}

	handlerMap := map[string]TravelHandler {
		"": defaultHandler,
	}

	options := travel.TravelOptions{
		StrictTraversal:   true,
		SubpathMaxLength: map[string]int{
			"GET":    travel.UnlimitedSubpath,
			"PUT":    0,
			"POST":   0,
			"DELETE": 0,
		},
	}
	r, err := travel.NewRouter(getRootTree, hm, errorHandler, &options)
	if err != nil {
		log.Fatalf("Error creating Travel router: %v\n", err)
	}
	http.Handle("/", r)
	http.ListenAndServe("127.0.0.1:8000", nil)


Travel provides additional options to modify normal traversal semantincs:

Strict vs. Permissive

"Strict" means to follow Pyramid traversal semantics -- handler name can only be "" (empty string) or the latest token in path when
root tree lookup failed (everything beyond that is the subpath). Note that this can be modified with handler name overrides in the
root tree object.

Non-strict (permissive) means that the handler name is always the latest token in the path (regardless if lookup fully succeeds).

Strict setting has no effect on the following options (they can be used to modify "strict" traversal as needed). Handler names
can always be overridden by embedding handler keys within the root tree ('%handler' key within the object, value must be a string).

Handler Overrides

Any level of the root tree can contain a special key "%handler", mapping to a handler name string that will be invoked instead of
whatever traversal would otherwise dictate. Handler overrides take precedence over strict/permissive mode rules.

Default Handler

The optional DefaultHandler is used to execute a fallback handler when traversal succeeds but the handler name returned is not
found within the handler map. Otherwise a 501 Not Implemented error is returned. As with handler overrides, the DefaultHandler
setting is respected regardless of strict/permissive setting.

Subpath Max Length

SubpathMaxLength is a map of method verb (all caps) to an integer representing the allowed number of subpath tokens. If the subpath
length is less than or equal to this limit, the request succeeds and the handler is executed per traversal semantics. If the subpath
exceeds this limit a 404 Not Found is returned. Traditional Pyramid traversal has an unlimited subpath max length. That can be emulated by setting SubpathMaxLength[verb] to
UnlimitedSubpath.

*/