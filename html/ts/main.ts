import { CardTemplate } from "./modules/card_template.js";
import { ChoiceTemplate } from "./modules/choice_template.js";
import { KeyValueTemplate } from "./modules/keyvalue_template.js";
import { ListItemTemplate } from "./modules/list_item_template.js";
import { WSocket } from "./websocket.js";

export let networkManager: WSocket;

const ready = (callback: () => void): void => {
    if (document.readyState !== 'loading') callback();
    else document.addEventListener('DOMContentLoaded', callback);
};

ready(() => {
    KeyValueTemplate.SetTemplate(document.getElementById('keyvalue-template') as HTMLTemplateElement);
    CardTemplate.SetTemplate(document.getElementById('card-template') as HTMLTemplateElement);
    CardTemplate.SetToLocation(document.querySelector('main') as HTMLElement);
    ChoiceTemplate.SetTemplate(document.getElementById('choice-template') as HTMLTemplateElement);
    ListItemTemplate.SetTemplate(document.getElementById('list-item-template') as HTMLTemplateElement);

    networkManager = new WSocket();
    (window as any).networkManager = networkManager;
});
