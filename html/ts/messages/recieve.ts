import { RecieveMessageType, Status } from "../enums";
export interface RecieveMessageEnvelope {
    readonly Type: RecieveMessageType,
    readonly ItemID: string,
    readonly Payload: string,
}
export interface CreateItem {
    readonly AddedAt: number,
    readonly Name: string,
    readonly Path: string,
    readonly Image: string,
    readonly Status: Status,
}
export interface AcceptToken {
    readonly Accept: boolean,
    readonly TokenID: string,
    readonly Name: string,
    readonly ImgSrc: string,
    readonly Data: Record<string, string>,
}
export interface AddMetadata {
    readonly KeyValues: Record<string, string>,
}
export interface ChangeStatus {
    readonly Status: Status,
}
export interface DetectingAddLogs {
    readonly Value: string,
    readonly Icon: string,
}
export interface Choice {
    readonly TokenID: string,
    readonly ImgSrc: string,
    readonly Data: Record<string, string>,
}
export interface AddChoices {
    readonly Choices: Choice[],
}
