package gonest

func AckMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_ACK,
		Body: &Event_AckMessage{
			AckMessage: &AckMessage{},
		},
		IpAddress: source,
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
		IpAddress: source,
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
		IpAddress: source,
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
		IpAddress: source,
	}
}

func HelloMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_HELLO,
		Body: &Event_HelloMessage{
			HelloMessage: &HelloMessage{},
		},
		IpAddress: source,
	}
}

func HumanVoteMessageFactory(source, target string) *Event {
	return &Event{
		MessageType: MessageType_HUMANVOTE,
		Body: &Event_HumanVoteMessage{
			HumanVoteMessage: &HumanVoteMessage{
				Target: target,
			},
		},
		IpAddress: source,
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
		IpAddress: source,
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
		IpAddress: source,
	}
}

func ItsHimMessageFactory(source string) *Event {
	return &Event{
		MessageType: MessageType_ITSHIM,
		Body: &Event_ItsHimMessage{
			ItsHimMessage: &ItsHimMessage{},
		},
		IpAddress: source,
	}
}

func WerewolfVoteMessageFactory(source, target string) *Event {
	return &Event{
		MessageType: MessageType_WEREWOLFVOTE,
		Body: &Event_WerwolfVoteMessage{
			WerwolfVoteMessage: &WerewolfVoteMessage{
				Target: target,
			},
		},
		IpAddress: source,
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
		IpAddress: source,
	}
}
