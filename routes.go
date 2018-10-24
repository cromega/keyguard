package main

func (s *server) routes() {
	s.router.Route("/").Get(s.rootHandler())
	s.router.Route("/:expiry").Get(s.rootHandler())
	s.router.Route("/key").Get(s.keyHandler())
	s.router.Route("/pubkey").Get(s.pubKeyHandler())
}
