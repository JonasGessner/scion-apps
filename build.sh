#!/usr/bin/env sh

# Deal with go stuff...
go mod download golang.org/x/mobile

mkdir .gomobilebuild >> /dev/null 2>&1

cd .gomobilebuild

function buildiOS {
    echo "Building ios (arm64)"
    GO11MODULE=on gomobile bind -target=ios/ios_arm64 -o _AppnetIOS.framework -iosversion 13.0 ../pkg/appnet && rm -rf AppnetIOS.framework && mv _AppnetIOS.framework AppnetIOS.framework && mv AppnetIOS.framework/Versions/A/_AppnetIOS AppnetIOS.framework/Versions/A/AppnetIOS.a
}

function buildSim {
    echo "Building simulator (x86_64)"
    GO11MODULE=on gomobile bind -target=ios/sim_amd64 -o _AppnetSim.framework -iosversion 13.0 ../pkg/appnet && rm -rf AppnetSim.framework && mv _AppnetSim.framework AppnetSim.framework && mv AppnetSim.framework/Versions/A/_AppnetSim AppnetSim.framework/Versions/A/AppnetSim.a
}

function buildCatalyst {
    echo "Building catalyst (x86_64)"
    GO11MODULE=on gomobile bind -target=ios/catalyst_amd64 -o _AppnetCatalyst.framework -macosversion 10.15 ../pkg/appnet && rm -rf AppnetCatalyst.framework && mv _AppnetCatalyst.framework AppnetCatalyst.framework && mv AppnetCatalyst.framework/Versions/A/_AppnetCatalyst AppnetCatalyst.framework/Versions/A/AppnetCatalyst.a
}


if [[ "$1" == "catalyst" ]]; then
    buildCatalyst
elif [[ "$1" == "ios" ]]; then
    buildiOS
elif  [[ "$1" == "sim" ]]; then
    buildSim
elif  [[ "$1" == "no-build" ]]; then
    echo "Skipping build"
elif [ "$#" -eq 0 ]; then
    echo "Building all"
    
    buildSim
    buildCatalyst
    buildiOS
else
    echo "Invalid arguments"
fi

cd ..

echo "Making xcframework"

xcodebuild -create-xcframework -library .gomobilebuild/AppnetSim.framework/Versions/A/AppnetSim.a -headers .gomobilebuild/AppnetSim.framework/Versions/A/Headers -library .gomobilebuild/AppnetIOS.framework/Versions/A/AppnetIOS.a -headers .gomobilebuild/AppnetIOS.framework/Versions/A/Headers -library .gomobilebuild/AppnetCatalyst.framework/Versions/A/AppnetCatalyst.a -headers .gomobilebuild/AppnetCatalyst.framework/Versions/A/Headers -output _Appnet.xcframework && rm -rf Appnet.xcframework && mv _Appnet.xcframework Appnet.xcframework

echo Done
