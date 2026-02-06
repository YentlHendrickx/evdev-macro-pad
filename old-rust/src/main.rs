use evdev::*;
#[allow(dead_code)]
#[derive(Debug, Copy, Clone)]
enum VolumeAction {
    Up,
    Down,
}

fn capture_input_events<F>(
    event_id: u16,
    mut on_volume: F,
) -> Result<(), Box<dyn std::error::Error>>
where
    F: FnMut(VolumeAction),
{
    let path = format!("/dev/input/event{}", event_id);
    let mut device = Device::open(path)?;

    loop {
        for event in device.fetch_events()? {
            match event.destructure() {
                EventSummary::Key(_, KeyCode::KEY_K, 1) => {
                    on_volume(VolumeAction::Up);
                }
                EventSummary::Key(_, KeyCode::KEY_J, 1) => {
                    on_volume(VolumeAction::Down);
                }
                _ => {}
            }
        }
    }
}

struct ZoomController {
    levels: Vec<i32>,
    index: i32,
}

impl ZoomController {
    fn new() -> Self {
        let zoom_min = 100;
        let zoom_max = 133;
        let count = 33;

        let levels = (0..count)
            .map(|i| zoom_min + i * (zoom_max - zoom_min) / (count - 1))
            .map(|v| v as i32)
            .collect();

        Self { levels, index: 0 }
    }

    fn set_zoom_level(level: i32) {
        use std::process::Command;
        let output = Command::new("v4l2-ctl")
            .arg("--set-ctrl")
            .arg(format!("zoom_absolute={}", level))
            .output()
            .expect("Failed to execute command");
        if output.status.success() {
            println!("Zoom set to level: {}", level);
        } else {
            eprintln!(
                "Failed to set zoom level: {}. Error: {}",
                level,
                String::from_utf8_lossy(&output.stderr)
            );
        }
    }

    fn handle(&mut self, action: VolumeAction) {
        match action {
            VolumeAction::Up => self.index += 1,
            VolumeAction::Down => self.index -= 1,
        }

        self.index = self.index.clamp(0, (self.levels.len() - 1) as i32);
        ZoomController::set_zoom_level(self.levels[self.index as usize]);
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut zoom = ZoomController::new();

    // TODO: Make event_id configurable -> device selection or just poll all devices?
    capture_input_events(7, |action| {
        zoom.handle(action);
    })?;

    Ok(())
}
