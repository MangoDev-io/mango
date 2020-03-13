import { Component, OnInit } from '@angular/core'
import { StateService } from '../../state.service'
import { AssetListing } from 'src/app/asset-listing'
@Component({
    selector: 'app-token-lister',
    templateUrl: './token-lister.component.html',
    styleUrls: ['./token-lister.component.scss'],
})
export class TokenListerComponent implements OnInit {
    ownedAssets: AssetListing[]

    constructor(private stateService: StateService) {}

    ngOnInit(): void {
        this.getOwnedAssets()
    }

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

    getOwnedAssets() {
        this.stateService.getAssetListings().subscribe(x => {
            this.ownedAssets = x
        })
    }
}
