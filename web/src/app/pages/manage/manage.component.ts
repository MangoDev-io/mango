import { Component, OnInit } from '@angular/core'
import { StateService } from 'src/app/state.service'
import { Subscription } from 'rxjs'
import { Token } from 'src/app/model/token'

@Component({
    selector: 'app-manage',
    templateUrl: './manage.component.html',
    styleUrls: ['./manage.component.scss'],
})
export class ManageComponent implements OnInit {
    showTokenSubscription: Subscription
    selectedTokenSubscription: Subscription

    showTokenCreate: boolean
    selectedToken: Token

    constructor(private stateService: StateService) {
        this.showTokenSubscription = this.stateService
            .getShowCreateToken()
            .subscribe(b => {
                this.showTokenCreate = b
            })

        this.selectedTokenSubscription = this.stateService
            .getSelectedToken()
            .subscribe(t => {
                this.selectedToken = t
            })
    }

    ngOnInit() {}
}
