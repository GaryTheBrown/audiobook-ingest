import { SendMessageType } from "../enums.js";
import { networkManager } from "../main.js";
import { CreateItem } from "../messages/recieve.js";
import { SubmitToken } from "../messages/send.js";
import { KeyValueTemplate } from "./keyvalue_template.js";
export class Header {
    private readonly itemID: string;
    private readonly addedAt: number;
    private name: string;
    private tokenSpan_DOM: HTMLElement;
    private img_DOM: HTMLImageElement;
    private name_DOM: HTMLElement;
    private path_DOM: HTMLElement;
    private tokenCode_DOM: HTMLElement;
    private tokenInput_DOM: HTMLInputElement;
    private tokenButton_DOM: HTMLButtonElement;
    private timeBadge_DOM: HTMLElement;
    private timeBadgeCode_DOM: HTMLElement;
    private metadataTable_DOM: HTMLTableElement;

    constructor(itemID: string, newItem: CreateItem, header_DOM: HTMLElement) {
        this.itemID = itemID;
        this.name = newItem.Title;
        this.addedAt = newItem.AddedAt;
        this.img_DOM = header_DOM.querySelector('img') as HTMLImageElement;
        this.img_DOM.src = newItem.Image;
        this.name_DOM = header_DOM.querySelector('h3') as HTMLElement;
        this.name_DOM.innerText = newItem.Title;
        this.path_DOM = header_DOM.querySelector('h4') as HTMLElement;
        this.path_DOM.innerText = newItem.Path;
        this.tokenSpan_DOM = header_DOM.querySelector('span') as HTMLElement;
        this.tokenCode_DOM = this.tokenSpan_DOM.querySelector('code') as HTMLElement;
        this.tokenInput_DOM = header_DOM.querySelector('input[name="token"]') as HTMLInputElement;
        this.tokenButton_DOM = header_DOM.querySelector('button.tokensubmit') as HTMLButtonElement;
        this.timeBadge_DOM = header_DOM.querySelector('h5') as HTMLElement;
        this.timeBadgeCode_DOM = this.timeBadge_DOM.querySelector('code') as HTMLElement;

        const dateObj = new Date(newItem.AddedAt * 1000);
        const formattedDate = dateObj.toLocaleDateString(undefined, {
            year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
        });

        this.timeBadgeCode_DOM.innerHTML = formattedDate;

        this.metadataTable_DOM = header_DOM.querySelector('details dialog table') as HTMLTableElement;
    }

    public AddEventListeners(): void {
        this.tokenButton_DOM.addEventListener("click", this.clickSubmitToken);
    }

    private clickSubmitToken = (e: Event): void => {
        this.tokenInput_DOM.hidden = true;
        this.tokenButton_DOM.hidden = true;
        const payloadStruct: SubmitToken = {
            TokenID: this.tokenInput_DOM.value,
        };
        const payload: string = JSON.stringify(payloadStruct);
        networkManager.SendPayload(this.itemID, SendMessageType.SubmitToken, payload);
    };

    public get AddedAt(): number { return this.addedAt; }
    public get Name(): string { return this.name; }
    public set Name(title: string) { this.name = title; this.name_DOM.innerText = title; }

    public AddMetadata(keyValues: KeyValueTemplate[]): void {
        this.RemoveMetadata();
        keyValues.forEach(keyValue => {
            this.metadataTable_DOM.appendChild(keyValue.Element() as HTMLTableRowElement);
        });
    }

    public RemoveMetadata(): void {
        this.metadataTable_DOM.innerHTML = "";
    }

    public AcceptSubmitToken(accept: Boolean, TokenID: string): void {
        if (accept) {
            this.tokenCode_DOM.innerText = TokenID;
            this.tokenSpan_DOM.hidden = false;
            this.tokenInput_DOM.hidden = true
            this.tokenButton_DOM.hidden = true
        } else {
            this.tokenInput_DOM.hidden = false;
            this.tokenButton_DOM.hidden = false;
        }
    }
}
