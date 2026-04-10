package discordgoplus

import (
	"fmt"

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

func ParseModalDataV2(
	modalData discordgo.ModalSubmitInteractionData,
) map[string]interface{} {
	data := make(map[string]interface{}, 0)
	for _, row := range modalData.Components {
		labelRow := row.(*discordgo.Label)
		cmpnt := labelRow.Component

		if cmpnt.Type() == 3 {
			input := cmpnt.(*discordgo.SelectMenu)

			if len(input.Values) == 0 {
				continue
			}

			if input.MaxValues == 1 || input.MaxValues == 0 {
				data[input.CustomID] = input.Values[0]
				continue
			}

			data[input.CustomID] = input.Values
		}

		if cmpnt.Type() == 4 {
			input := cmpnt.(*discordgo.TextInput)
			data[input.CustomID] = input.Value
			continue
		}

		if cmpnt.Type() == 19 {
			input := cmpnt.(*discordgo.FileUpload)

			if input.MaxValues == 1 || input.MaxValues == 0 {
				if len(input.Values) == 0 {
					continue
				}

				data[input.CustomID] = modalData.Resolved.Attachments[input.Values[0]]
				continue
			}

			values := make([]*discordgo.MessageAttachment, 0)
			for _, value := range input.Values {
				values = append(values, modalData.Resolved.Attachments[value])
			}

			data[input.CustomID] = values
		}

		fmt.Println(cmpnt.Type())
	}

	return data
}
