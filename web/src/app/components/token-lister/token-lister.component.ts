import { Component, OnInit } from '@angular/core'
import { StateService } from '../../state.service'

@Component({
    selector: 'app-token-lister',
    templateUrl: './token-lister.component.html',
    styleUrls: ['./token-lister.component.scss'],
})
export class TokenListerComponent implements OnInit {
    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    createNewAsset() {
        this.stateService.setShowCreateToken(true)
        let entries = document.getElementsByClassName('token-entry__container')
        for (let i = 0; i < entries.length; i++) {
            entries[i].classList.remove('active')
        }

        document
            .getElementsByClassName('create-new__container')[0]
            .classList.add('active')
    }
}
