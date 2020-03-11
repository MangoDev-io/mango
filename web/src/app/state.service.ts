import { Injectable } from '@angular/core'
import { Token } from './token'
import { AssetRequest } from './assetRequest'
import { BehaviorSubject, Observable, from } from 'rxjs'
import { HttpClient, HttpHeaders } from '@angular/common/http'
import { environment } from '../environments/environment'
import algosdk from 'algosdk'
import { AssetListing } from './asset-listing'

@Injectable({
    providedIn: 'root',
})
export class StateService {
    //private baseURL = 'http://localhost:5000'
    private baseURL = 'http://api.mangodev.io'

    private authToken: string

    private selectedTokenSubject = new BehaviorSubject<Token>(null)
    private showCreateSubject = new BehaviorSubject<boolean>(false)

    private algorandClient = new algosdk.Algod(
        environment.algorandToken,
        environment.algorandAddress,
        ''
    )

    constructor(private httpClient: HttpClient) {
        this.showCreateSubject.next(true)
    }

    getSelectedToken(): Observable<Token> {
        return this.selectedTokenSubject.asObservable()
    }

    setSelectedToken(token: Token) {
        this.selectedTokenSubject.next(token)
    }

    getShowCreateToken(): Observable<boolean> {
        return this.showCreateSubject.asObservable()
    }

    setShowCreateToken(b: boolean) {
        this.showCreateSubject.next(b)
    }

    setAuthToken(a: string) {
        this.authToken = a
    }

    encodeMnemonic(m: string) {
        return this.httpClient.post(this.baseURL + '/encodeMnemonic', {
            mnemonic: m,
        })
    }

    getAssetListings(): Observable<AssetListing[]> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.get<AssetListing[]>(
            this.baseURL + '/assets',
            options
        )
    }

    createAsset(a: Token) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/createAsset', a, options)
    }

    freezeAsset(a: AssetRequest) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/freezeAsset', a, options)
    }

    revokeAsset(a: AssetRequest) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/revokeAsset', a, options)
    }

    modifyAsset(a: AssetRequest) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/modifyAsset', a, options)
    }

    destroyAsset(a: AssetRequest) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/destroyAsset', a, options)
    }

    getAssetDetails(assetId: string) {
        return this.algorandClient.assetInformation(assetId)
    }
}
