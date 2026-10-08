import { RecieveMessageType, SendMessageType } from "./enums";
import { AcceptToken, AddChoices, ChangeStatus, CreateItem, DetectingAddLogs, RecieveMessageEnvelope } from "./messages/recieve";
import { SendMessageEnvelope } from "./messages/send";
import { CardTemplate } from "./modules/card_template";
import { ChoiceTemplate } from "./modules/choice_template";
import { KeyValueTemplate } from "./modules/keyvalue_template";
import { ListItemTemplate } from "./modules/list_item_template";
export class WSocket {
    private wsProtocol: string;
    private webSocket: WebSocket;
    private cardRegistry = new Map<string, CardTemplate>();
    private statusIndicator_DOM: HTMLElement;

    private pingIntervalID: number | null = null;

    constructor() {
        this.wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
        this.statusIndicator_DOM = document.querySelector('header span[data-connection]') as HTMLElement;
        this.setIndicatorState("connecting");
        this.webSocket = new WebSocket(this.wsProtocol + window.location.host + '/ws');
        this.pairFunctions();
    }

    private connect(): void {
        this.setIndicatorState("connecting");
        this.webSocket = new WebSocket(this.wsProtocol + window.location.host + '/ws');
        this.pairFunctions();
    }
    private pairFunctions(): void {
        this.webSocket.onopen = this.onOpen;
        this.webSocket.onmessage = this.onMessage;
        this.webSocket.onclose = this.onClose;
        this.webSocket.onerror = this.onError;
    }

    private setIndicatorState(state: "connecting" | "connected" | "disconnected"): void {
        if (this.statusIndicator_DOM) {
            this.statusIndicator_DOM.dataset.connection = state;
        }
    }

    public SendPayload(itemID: string, type: SendMessageType, payload: string = ""): void {
        if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
            const envelope: SendMessageEnvelope = {
                ItemID: itemID,
                Type: type,
                Payload: payload,
            };
            this.webSocket.send(JSON.stringify(envelope));
        }
    }

    private onOpen = (): void => {
        this.setIndicatorState("connected");

        this.pingIntervalID = window.setInterval(() => {
            this.SendPayload("SYSTEM", SendMessageType.Ping, `"PING"`);
        }, 30000);

        window.setTimeout(() => {
            const stuckIDs: string[] = [];
            for (const [itemID, card] of this.cardRegistry.entries()) {
                if (card.Status === "detecting") {
                    stuckIDs.push(itemID);
                }
            }

            if (stuckIDs.length > 0) {
                const serializedPayload = JSON.stringify(stuckIDs);
                console.log("stuckIDs: ", stuckIDs);
                this.SendPayload("SYSTEM", SendMessageType.DetectingStuck, serializedPayload);
            }
        }, 500);
    };

    private onClose = (): void => {
        this.setIndicatorState("disconnected");

        if (this.pingIntervalID) {
            window.clearInterval(this.pingIntervalID);
            this.pingIntervalID = null;
        }

        setTimeout(() => this.connect(), 3000);
    };

    private onError = (error: Event): void => {
        this.setIndicatorState("disconnected");
    };

    private onMessage = (e: MessageEvent): void => {
        try {
            const envelope = JSON.parse(e.data) as RecieveMessageEnvelope;
            const handler = this.messageHandlers[envelope.Type];
            if (handler) handler(envelope.ItemID, envelope.Payload);
        } catch (error) {
        }
    };

    private messageHandlers: Record<RecieveMessageType, (itemID: string, payload: string) => void> = {
        [RecieveMessageType.CreateItem]: (itemID: string, payload: string): void => {
            const itemData = JSON.parse(payload) as CreateItem;
            if (this.cardRegistry.has(itemID)) {
                console.error(`[SOCKET] Card ID ${itemID} already active on dashboard. Diverting to status update.`);
                const existingItem = this.cardRegistry.get(itemID);
                if (existingItem) {
                    existingItem.Status = itemData.Status;
                }
                return;
            }

            const item = new CardTemplate(itemID, itemData) as CardTemplate;
            item.InsertToDom();
            this.cardRegistry.set(itemID, item);
        },
        [RecieveMessageType.RemoveItem]: (itemID: string, payload: string): void => {
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) {
                item.Remove();
                this.cardRegistry.delete(itemID);
            }
        },
        [RecieveMessageType.AcceptToken]: (itemID: string, payload: string): void => {
            const acceptData = JSON.parse(payload) as AcceptToken;
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) {
                item.Header.AcceptSubmitToken(acceptData.Accept, acceptData.TokenID);
                if (acceptData.Accept) {
                    item.Header.Name = acceptData.Name
                    // TODO HERE IS WHERE WE ARE GOING TO ADD IN THE ITEM DATA IN CONVERTING SECTION FROM AcceptToken
                    //item.SectionConverting.

                }
            }

        },
        [RecieveMessageType.ChangeStatus]: (itemID: string, payload: string): void => {
            const statusData = JSON.parse(payload) as ChangeStatus;
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) item.Status = statusData.ToStatus;
        },
        [RecieveMessageType.AddMetadata]: (itemID: string, payload: string): void => {
            const rawFieldsData = JSON.parse(payload) as Record<string, string>;
            const KeyValueArray: KeyValueTemplate[] = [];
            for (const [key, value] of Object.entries(rawFieldsData)) {
                KeyValueArray.push(new KeyValueTemplate(key, value));
            }
            const item = this.cardRegistry.get(itemID);
            if (item) item.AddMetadata(KeyValueArray);
        },
        [RecieveMessageType.RemoveMetadata]: (itemID: string, payload: string): void => {
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) item.RemoveMetadata();
        },
        [RecieveMessageType.DetectingAddLog]: (itemID: string, payload: string): void => {
            const addListData = JSON.parse(payload) as DetectingAddLogs;
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            const tmpListItem = new ListItemTemplate(addListData.Value, addListData.Icon) as ListItemTemplate;
            if (item) item.SectionDetecting.AddLog(tmpListItem);
        },
        [RecieveMessageType.DetectingRemoveLogs]: (itemID: string, payload: string): void => {
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) item.SectionDetecting.RemoveLogs();
        },
        [RecieveMessageType.AddChoices]: (itemID: string, payload: string): void => {
            const item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) {
                const addChoicesData = JSON.parse(payload) as AddChoices;
                const finalChoices: ChoiceTemplate[] = [];
                addChoicesData.Choices.forEach(choiceIn => {
                    const keyvalueData: KeyValueTemplate[] = [];
                    for (const [key, value] of Object.entries(choiceIn.Data)) {
                        keyvalueData.push(new KeyValueTemplate(key, value));
                    }
                    finalChoices.push(new ChoiceTemplate(choiceIn.TokenID, choiceIn.ImgSrc, keyvalueData));
                });
                item.SectionMultichoice.AddChoices(finalChoices);
            }
        },
        [RecieveMessageType.RemoveChoices]: (itemID: string, payload: string): void => {
            let item = this.cardRegistry.get(itemID) as CardTemplate;
            if (item) item.SectionMultichoice.RemoveChoices();
        },
    };

}