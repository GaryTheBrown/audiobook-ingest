import { SendMessageType } from "../enums";
export interface SendMessageEnvelope {
    readonly Type: SendMessageType,
    readonly ItemID: string,
    readonly Payload: string,
}
export interface SubmitToken {
    readonly TokenID: string,
}
export interface SelectChoice {
    readonly Choice: string,
}