import { Component, OnInit } from '@angular/core'
import { StateService } from 'src/app/state.service'

@Component({
    selector: 'app-manage',
    templateUrl: './manage.component.html',
    styleUrls: ['./manage.component.scss'],
})
export class ManageComponent implements OnInit {
    public showTokenCreate = true

    public mnemonic: string

    constructor(private stateService: StateService) {}

    ngOnInit() {
        this.mnemonic = this.stateService.getMnemonic()
    }
}
