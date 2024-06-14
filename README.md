# caprice

-  This will be software which will allow you to change your robot's software (therefore, personality) on the fly.

## Dev notes

### How will I architect this?

-  Write in Go, both client (Vector) and server
-  Server probably won't need much, just a JSON of available firmwares and be able to act as a fileserver
-  I will need a way to pack up a firmware. This will probably need to happen manually.
-  Client will run as a seperate daemon (multi-user.target). There will be an install script (or installer program?)
-  I will need a patch system.
  -  Patches:
      -  ln /dev/ttyHS0 /dev/ttyHSL1 (at bootup as well. could just implement this as part of a custom ankiinit regardless of firmware)
      -  increase CPU and RAM frequencies? some software runs a little too slow on the default (also part of ankiinit)
      -  mm-anki-camera/mm-qcamera-daemon swap
          -  There seems to be three eras of mm programs.
            -  Feb 9 2018 - DVT2
            -  0.10 - DVT3
            -  0.14(?) and above - modern robots
            -  Android - i refuse to deal with Android-era software. too much needs to change
      -  Copy token.jwt from com.anki.victor and change server endpoints
          -  probably only for 0.10
          -  0.10 also needs: `sed -i "s/robot_id/token_id/g" /anki/bin/vic-cloud`
   -   So maybe like:
```go
const (
  CameraEra_DVT2 = 0
  CameraEra_DVT3 = 1
  CameraEra_Modern = 2
)

type PatchConfig struct {
  CameraEra int
  DoCloudPatches bool
}
```
  -  That's theoretically it... everything else could be handled by a custom ankiinit and won't hurt any modern software
