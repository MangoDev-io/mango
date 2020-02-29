export class Token {
    assetId: string
    creatorAddr: string
    assetName: string
    unitName: string
    total: number
    decimals: number
    defaultFrozen: boolean
    url: string
    metadataHash: string
    managerAddr: string
    reserveAddr: string
    freezeAddr: string
    clawbackAddr: string

    constructor(
        assetId: string,
        creatorAddr: string,
        assetName: string,
        unitName: string,
        total: number,
        decimals: number,
        defaultFrozen: boolean,
        url: string,
        metadataHash: string,
        managerAddr: string,
        reserveAddr: string,
        freezeAddr: string,
        clawbackAddr: string
    ) {
        this.assetId = assetId
        this.creatorAddr = creatorAddr
        this.assetName = assetName
        this.unitName = unitName
        this.total = total
        this.decimals = decimals
        this.defaultFrozen = defaultFrozen
        this.url = url
        this.metadataHash = metadataHash
        this.managerAddr = managerAddr
        this.reserveAddr = reserveAddr
        this.freezeAddr = freezeAddr
        this.clawbackAddr = clawbackAddr
    }
}
