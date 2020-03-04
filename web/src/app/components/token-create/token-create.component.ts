import { Component, OnInit } from '@angular/core'
import { AssetCreate } from 'src/app/asset-create'
import { StateService } from 'src/app/state.service'

@Component({
    selector: 'app-token-create',
    templateUrl: './token-create.component.html',
    styleUrls: ['./token-create.component.scss'],
})
export class TokenCreateComponent implements OnInit {
    private assetCreate = new AssetCreate()

    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    createAsset() {
        console.log(JSON.stringify(this.assetCreate))
        this.stateService.createAsset(this.assetCreate).subscribe(x => {
            console.log(x)
        })
    }
}
