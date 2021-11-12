package ipamlogic

/*
func updateCache(log logging.Logger, c *cache.Cache, prefix *gnmi.Path, path *gnmi.Path, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	//log.Debug("Debug updateState", "refPaths", refPaths)
	log.Debug("Debug updateState", "data", x)
	n, err := c.GetNotificationFromJSON(prefix, path, x, refPaths)
	if err != nil {
		return err
	}

	//printNotification(log, n)
	if n != nil {
		if err := c.GnmiUpdate(prefix.Target, n); err != nil {
			if strings.Contains(fmt.Sprintf("%v", err), "stale") {
				return nil
			}
			return err
		}
	}
	return nil
}

func deleteCache(log logging.Logger, c *cache.Cache, prefix *gnmi.Path, path *gnmi.Path) error {
	n, err := c.GetNotificationFromDelete(prefix, path)
	if err != nil {
		return err
	}
	if err := c.GnmiUpdate(prefix.Target, n); err != nil {
		return err
	}

	return nil
}
*/

/*
func printNotification(log logging.Logger, n *gnmi.Notification) {
	log.Debug("Debug Notification", "Notification", n)
	for _, u := range n.GetUpdate() {
		log.Debug("Update", "Path", u.GetPath(), "Value", u.GetVal())
	}
}
*/

/*
func getCache(log logging.Logger, c *cache.Cache, prefix, p *gnmi.Path) (*gnmi.TypedValue, error) {
	n, err := c.Query(prefix.Target, prefix, p)
	if err != nil {
		return nil, err
	}
	log.Debug("Get Cache", "Path", p, "Notification", n)
	if n != nil {
		if len(n.GetUpdate()) > 0 {
			return n.GetUpdate()[0].GetVal(), nil
		}
	}
	return nil, nil
}
*/
