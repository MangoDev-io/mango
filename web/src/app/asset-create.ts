export class AssetCreate {
    creatorAddr: string
    assetName: string
    unitName: string
    totalIssuance: number
    decimals: number
    defaultFrozen: boolean
    url: string
    metadataHash: string
    managerAddr: string
    reserveAddr: string
    freezeAddr: string
    clawbackAddr: string

    // Assigns parameters from object to class fields
    constructor(init?: Partial<AssetCreate>) {
        Object.assign(this, init)
    }
}
