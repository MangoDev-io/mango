import { Component, OnInit } from '@angular/core'
import { StateService } from 'src/app/state.service'
import { Token } from '../../token'

@Component({
    selector: 'app-token-create',
    templateUrl: './token-create.component.html',
    styleUrls: ['./token-create.component.scss'],
})
export class TokenCreateComponent implements OnInit {
    assetCreate = new Token()

    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    createAsset() {
        console.log(JSON.stringify(this.assetCreate))
        this.stateService.createAsset(this.assetCreate).subscribe(x => {
            console.log(x)
        })
    }
}
