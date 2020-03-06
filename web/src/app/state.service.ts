import { Injectable } from '@angular/core'
import { Token } from './model/token'
import { BehaviorSubject, Observable } from 'rxjs'
import {
    HttpClient,
    HttpHeaderResponse,
    HttpHeaders,
} from '@angular/common/http'
import { AssetCreate } from './asset-create'

@Injectable({
    providedIn: 'root',
})
export class StateService {
    private baseURL = 'http://localhost:5000'

    public currToken = new Token(
        '324234',
        'MTTIYHO4QQAEZ6VPLW44YV7AT7TY4SDBYGLFARD7POC2JMSIIK3RVWZVZU',
        'USD Tether',
        'USDT',
        100000000,
        15,
        false,
        'https://usdtether.io',
        'aldskfjlakjsdf;',
        'MTTIYHO4QQAEZ6VPLW44YV7AT7TY4SDBYGLFARD7POC2JMSIIK3RVWZVZU',
        'MTTIYHO4QQAEZ6VPLW44YV7AT7TY4SDBYGLFARD7POC2JMSIIK3RVWZVZU',
        'MTTIYHO4QQAEZ6VPLW44YV7AT7TY4SDBYGLFARD7POC2JMSIIK3RVWZVZU',
        'MTTIYHO4QQAEZ6VPLW44YV7AT7TY4SDBYGLFARD7POC2JMSIIK3RVWZVZU'
    )

    public tokenList: Token[] = [this.currToken, this.currToken, this.currToken]

    private authToken: string

    private selectedTokenSubject = new BehaviorSubject<Token>(null)
    private showCreateSubject = new BehaviorSubject<boolean>(false)

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

    createAsset(a: AssetCreate) {
        let httpHeaders = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + this.authToken,
        })
        let options = {
            headers: httpHeaders,
        }

        return this.httpClient.post(this.baseURL + '/createAsset', a, options)
    }
}
