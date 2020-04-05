package gonest

func AckMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_ACK,
		Body: &Event_AckMessage{
			AckMessage: &AckMessage{},
		},
		Source: source,
	}
}

func ChatMessageFactory(source, message string) *Event {
	return &Event{
		MessageType: MessageType_CHAT,
		Body: &Event_ChatMessage{
			ChatMessage: &ChatMessage{
				Content: message,
			},
		},
		Source: source,
	}
}

func ClairvoyantMessageFactory(source, target string) *Event {
	return &Event{
		MessageType: MessageType_CLAIRVOYANT,
		Body: &Event_ClairvoyantMessage{
			ClairvoyantMessage: &ClairvoyantMessage{
				Target: target,
			},
		},
		Source: source,
	}
}

func CupidMessageFactory(source, target1, target2 string) *Event {
	return &Event{
		MessageType: MessageType_CUPID,
		Body: &Event_CupidMessage{
			CupidMessage: &CupidMessage{
				IpAddress1: target1,
				IpAddress2: target2,
			},
		},
		Source: source,
	}
}

func DeadMessageFactory(source string, role Role, reason Reason) *Event {
	return &Event{
		MessageType: MessageType_DEAD,
		Body: &Event_DeadMessage{
			DeadMessage: &DeadMessage{
				Role:   role,
				Reason: reason,
			},
		},
		Source: source,
	}
}

func HelloMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_HELLO,
		Body: &Event_HelloMessage{
			HelloMessage: &HelloMessage{},
		},
		Source: source,
	}
}

func HunterMessageFactory(source, target string) *Event {
	return &Event{
		MessageType: MessageType_HUNTER,
		Body: &Event_HunterMessage{
			HunterMessage: &HunterMessage{
				Target: target,
			},
		},
		Source: source,
	}
}

func IpListMessageFactory(source, ipAdress string) *Event {
	ipAdressList := make([]string, 1)
	ipAdressList[0] = ipAdress
	return &Event{
		MessageType: MessageType_IPLIST,
		Body: &Event_IpListMessage{
			IpListMessage: &IpListMessage{
				IpAdress: ipAdressList,
			},
		},
		Source: source,
	}
}

func ItsHimMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_ITSHIM,
		Body: &Event_ItsHimMessage{
			ItsHimMessage: &ItsHimMessage{},
		},
		Source: source,
	}
}

func LeaderElectionMessageFactory(source, leader string) *Event {
	return &Event{
		MessageType: MessageType_LEADERELECTION,
		Body: &Event_LeaderElectionMessage{
			LeaderElectionMessage: &LeaderElectionMessage{
				Leader: leader,
			},
		},
		Source: source,
	}
}

func RoleDistributionMessageFactory(source, target string, role Role) *Event {
	return &Event{
		MessageType: MessageType_ROLEDISTRIBUTION,
		Body: &Event_RoleDistributionMessage{
			RoleDistributionMessage: &RoleDistributionMessage{
				Target: target,
				Role:   role,
			},
		},
		Source: source,
	}
}

func VoteMessageFactory(source, target string) *Event {
	return &Event{
		MessageType: MessageType_VOTE,
		Body: &Event_VoteMessage{
			VoteMessage: &VoteMessage{
				Target: target,
			},
		},
		Source: source,
	}
}

func WitchMessageFactory(source, target string, action WitchAction) *Event {
	return &Event{
		MessageType: MessageType_WITCH,
		Body: &Event_WitchMessage{
			WitchMessage: &WitchMessage{
				Action: action,
				Target: target,
			},
		},
		Source: source,
	}
}
