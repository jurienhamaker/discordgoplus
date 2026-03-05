package discordgoplus

import (
	"github.com/bwmarrin/discordgo"
)

func ParseModalData(
	modalData discordgo.ModalSubmitInteractionData,
) map[string]string {
	data := make(map[string]string, 0)
	for _, row := range modalData.Components {
		actionRow := row.(*discordgo.ActionsRow)

		for _, cmpnt := range actionRow.Components {
			if cmpnt.Type() == 4 {
				txtInput := cmpnt.(*discordgo.TextInput)
				data[txtInput.CustomID] = txtInput.Value
			}
		}
	}

	return data
}
