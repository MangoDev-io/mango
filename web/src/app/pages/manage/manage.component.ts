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
    private showTokenSubscription: Subscription
    private selectedTokenSubscription: Subscription

    private showTokenCreate: boolean
    private selectedToken: Token

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
