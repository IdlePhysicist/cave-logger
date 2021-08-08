package gui

func (g *Gui) startMonitoring() {
	stop := make(chan int, 1)
	g.state.stopChans["trips"] = stop
	g.state.stopChans["caves"] = stop
	g.state.stopChans["cavers"] = stop
	go g.tripsPanel().monitoringTrips(g)
	go g.cavesPanel().monitoringCaves(g)
	go g.caversPanel().monitoringCavers(g)
}

func (g *Gui) stopMonitoring() {
	g.state.stopChans["trips"] <- 1
	g.state.stopChans["caves"] <- 1
	g.state.stopChans["cavers"] <- 1
}

/*func (g *Gui) updateTask() {
	g.app.QueueUpdateDraw(func() {
		g.caversPanel().setEntries(g) // REVIEW: Why is this just people ?
	})
}*/
