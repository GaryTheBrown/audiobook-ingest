import { Status, StatusFromInt } from "../enums.js"
import { Header } from "./card_template_header.js";
import { SectionDetecting } from "../stages/sections_detecting.js";
import { SectionMultichoice } from "../stages/sections_multichoice.js";
import { SectionManual } from "../stages/sections_manual.js";
import { SectionConverting } from "../stages/sections_converting.js";
import { SectionVerifying } from "../stages/sections_verifying.js";
import { SectionVerified } from "../stages/sections_verified.js";
import { SectionFailed } from "../stages/sections_failed.js";
import { CreateItem } from "../messages/recieve.js";
import { KeyValueTemplate } from "./keyvalue_template.js";
export class CardTemplate {
    private static Template_DOM: HTMLTemplateElement;
    private static ToLocation_DOM: HTMLElement;

    public static SetTemplate(template: HTMLTemplateElement): void { CardTemplate.Template_DOM = template; }
    public static SetToLocation(location: HTMLElement): void { CardTemplate.ToLocation_DOM = location; }

    private readonly itemID: string;
    private readonly addedAt: number;

    private status: Status;
    private articleNode_DOM: HTMLElement;
    public Header: Header;
    public SectionDetecting: SectionDetecting;
    public SectionMultichoice: SectionMultichoice;
    public SectionManual: SectionManual;
    public SectionConverting: SectionConverting;
    public SectionVerifying: SectionVerifying;
    public SectionVerified: SectionVerified;
    public SectionFailed: SectionFailed;

    constructor(itemID: string, newItem: CreateItem) {
        this.articleNode_DOM = document.importNode(CardTemplate.Template_DOM.content, true).querySelector('article') as HTMLElement;
        this.itemID = itemID;
        this.status = newItem.Status;
        this.addedAt = newItem.AddedAt;
        this.articleNode_DOM.dataset.status = StatusFromInt[newItem.Status];
        this.Header = new Header(itemID, newItem, this.articleNode_DOM.querySelector('header') as HTMLElement);
        this.SectionDetecting = new SectionDetecting(this.articleNode_DOM.querySelector('section.detecting') as HTMLElement);
        this.SectionMultichoice = new SectionMultichoice(this.itemID, this.articleNode_DOM.querySelector('section.multichoice') as HTMLElement);
        this.SectionManual = new SectionManual(this.itemID, this.articleNode_DOM.querySelector('section.manual') as HTMLElement);
        this.SectionConverting = new SectionConverting(this.articleNode_DOM.querySelector('section.converting') as HTMLElement);
        this.SectionVerifying = new SectionVerifying(this.itemID, this.articleNode_DOM.querySelector('section.verifying') as HTMLElement);
        this.SectionVerified = new SectionVerified(this.itemID, this.articleNode_DOM.querySelector('section.verified') as HTMLElement);
        this.SectionFailed = new SectionFailed(this.itemID, this.articleNode_DOM.querySelector('section.failed') as HTMLElement);
    }

    private AddEventListeners(): void {
        this.Header.AddEventListeners();
        this.SectionMultichoice.AddEventListeners();
        this.SectionManual.AddEventListeners();
        this.SectionVerifying.AddEventListeners();
        this.SectionVerified.AddEventListeners();
        this.SectionFailed.AddEventListeners();
    }

    public Remove(): void {
        this.articleNode_DOM.remove();
    }

    public InsertToDom(skipListner: boolean = false): void {
        CardTemplate.ToLocation_DOM.appendChild(this.articleNode_DOM);
        if (!skipListner) {
            this.AddEventListeners();
        }
    }

    public get AddedAt(): number {
        return this.addedAt;
    }

    public get Status(): string {
        return StatusFromInt[this.status];
    }

    public set Status(status: Status) {
        this.articleNode_DOM.dataset.status = StatusFromInt[status];
        this.status = status;
    }

    public get ItemID(): string {
        return this.itemID;
    }

    public AddMetadata(keyValues: KeyValueTemplate[]): void {
        this.Header.AddMetadata(keyValues);
        this.SectionMultichoice.AddMetadata(keyValues);
    }

    public RemoveMetadata(): void {
        this.Header.RemoveMetadata();
        this.SectionMultichoice.RemoveMetadata();
    }
}
