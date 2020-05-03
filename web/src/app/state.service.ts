import { Injectable } from '@angular/core'
import { Token } from './token'
import { AssetRequest } from './assetRequest'
import { BehaviorSubject, Observable } from 'rxjs'
import { HttpClient, HttpHeaders } from '@angular/common/http'
import { AssetListing } from './asset-listing'
import { environment } from '../environments/environment'
import algosdk from 'algosdk'
import { Response } from './response'

@Injectable({
    providedIn: 'root',
})
export class StateService {
    private baseURL = 'http://localhost:5000'
    // private baseURL = 'https://api.mangodev.io'

    private authToken: string
    private address: string

    private selectedTokenSubject = new BehaviorSubject<Token>(null)
    private showCreateSubject = new BehaviorSubject<boolean>(false)
    private reloadListingsSubject = new BehaviorSubject<void>(null)

    private testnetAlgorandClient = new algosdk.Algod(
        environment.algorandToken,
        environment.testnetAlgorandAddress,
        ''
    )

    private mainnetAlgorandClient = new algosdk.Algod(
        environment.algorandToken,
        environment.mainnetAlgorandAddress,
        ''
    )

    public activeNetwork: string = 'testnet'

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

    setAddress(a: string) {
        this.address = a
    }

    getAddress(): string {
        return this.address
    }

    setReloadListings() {
        this.reloadListingsSubject.next()
    }

    getReloadListings() {
        return this.reloadListingsSubject.asObservable()
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
            this.baseURL + '/assets/' + this.activeNetwork,
            options
        )
    }

    createAsset(a: Token): Observable<Response> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post<Response>(
            this.baseURL + '/createAsset/' + this.activeNetwork,
            a,
            options
        )
    }

    freezeAsset(a: AssetRequest): Observable<Response> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post<Response>(
            this.baseURL + '/freezeAsset/' + this.activeNetwork,
            a,
            options
        )
    }

    revokeAsset(a: AssetRequest): Observable<Response> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post<Response>(
            this.baseURL + '/revokeAsset/' + this.activeNetwork,
            a,
            options
        )
    }

    modifyAsset(a: AssetRequest): Observable<Response> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post<Response>(
            this.baseURL + '/modifyAsset/' + this.activeNetwork,
            a,
            options
        )
    }

    destroyAsset(a: AssetRequest): Observable<Response> {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post<Response>(
            this.baseURL + '/destroyAsset/' + this.activeNetwork,
            a,
            options
        )
    }

    getAssetDetails(assetId: string) {
        if (this.activeNetwork == 'testnet') {
            return this.testnetAlgorandClient.assetInformation(assetId)
        } else if (this.activeNetwork == 'mainnet') {
            return this.mainnetAlgorandClient.assetInformation(assetId)
        }
    }
}
